package gateway

import (
	"fmt"
	"github.com/bufbuild/protovalidate-go"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/proto"
)

// ValidateHandler validates the given input message.
//
// It takes a proto.Message as input and returns an error.
func ValidateHandler(in proto.Message) error {
	validator, err := protovalidate.New()
	if err != nil {
		slog.Error(fmt.Sprintf("create user validator error: %v", err))
		return err
	}
	if err := validator.Validate(in); err != nil {
		return err
	}
	return nil
}
