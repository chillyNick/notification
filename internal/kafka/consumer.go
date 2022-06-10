package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/homework3/notification/internal/config"
	"github.com/homework3/notification/internal/notification"
	"github.com/homework3/notification/internal/repository"
	"github.com/homework3/notification/internal/stmp_sender"
	"github.com/homework3/notification/internal/tracer"
	"github.com/homework3/notification/pkg/model"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type consumerHandler struct {
	repo   repository.Repository
	sender *stmp_sender.MailSender
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *consumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Info().Msg("Setup consumer group session")

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Info().Msg("cleanup")

	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Info().Msg(fmt.Sprintf("Start consumer loop for topic: %s", claim.Topic()))
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// <https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29>
	for msg := range claim.Messages() {
		log.Info().
			Str("value", string(msg.Value)).
			Msgf("Message topic:%q partition:%d offset:%d", msg.Topic, msg.Partition, msg.Offset)

		spanCtx, err := tracer.ExtractSpanContext(msg.Headers)
		if err != nil {
			log.Error().Err(err).Msg("Failed to extract spanContext from kafka consumer headers")
		}
		span := opentracing.StartSpan("Comment after moderation processing", opentracing.ChildOf(spanCtx))

		comment := model.ModerationComment{}
		err = json.Unmarshal(msg.Value, &comment)
		if err != nil {
			log.Error().Err(err).Str("value", string(msg.Value)).Msg("Failed to unmarshal comment")
			span.Finish()

			continue
		}

		email, err := c.repo.GetMail(session.Context(), comment.UserId)
		if err != nil {
			log.Error().Err(err).Msg("Failed to get email address")
			span.Finish()

			continue
		}

		err = c.sender.SendMail(email, "Comment", notification.GetMailMessage(comment))
		if err != nil {
			log.Error().Err(err).Msg("Failed to send notification")
			//TODO add pushing into retry topic and add logic about processing it
		}

		session.MarkMessage(msg, "")
		span.Finish()
	}

	return nil
}

func StartProcessMessages(ctx context.Context, repo repository.Repository, sender *stmp_sender.MailSender, cfg *config.Kafka) error {
	consumer, err := createConsumerGroup(cfg)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create a consumer group")

		return err
	}
	defer consumer.Close()

	handler := &consumerHandler{repo: repo, sender: sender}
	loop := true
	for loop {
		err = consumer.Consume(ctx, []string{cfg.ConsumerTopic}, handler)
		if err != nil {
			log.Error().Err(err).Msg(" Consumer group session error")
		}

		select {
		case <-ctx.Done():
			loop = false
		default:

		}
	}

	return nil
}

func createConsumerGroup(cfg *config.Kafka) (sarama.ConsumerGroup, error) {
	saramaCfg := sarama.NewConfig()
	saramaCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	return sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupId, saramaCfg)
}
