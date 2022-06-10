package notification

import (
	"fmt"
	"github.com/homework3/notification/pkg/model"
	"github.com/rs/zerolog/log"
)

func GetMailMessage(comment model.ModerationComment) string {
	switch comment.Status {
	case model.ModerationCommentStatusPassed:
		return fmt.Sprintf("Your message: '%d' was published", comment.CommentId)
	case model.ModerationCommentStatusFailed:
		return fmt.Sprintf("Your message: '%d' was rejected by reason: '%s'", comment.CommentId, comment.Reason)
	}

	log.Warn().Msgf("unknown status %s", comment.Status)

	return fmt.Sprintf("Your message: '%s' was processing with status: '%s'", comment.CommentId, comment.Status)
}
