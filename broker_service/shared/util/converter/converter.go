package converter

import (
	"fmt"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func ConvertXByte2UUID(uuid [16]byte) string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

// targetType 必須是已定義proto
//
// 返回的是proto.Message，還要再使用斷言轉換成你想要的type
func ConvertAnyToType(anyMessage *anypb.Any, targetType reflect.Type) (proto.Message, error) {
	// 使用反射來創建一個目標類型的實例
	instance := reflect.New(targetType.Elem()).Interface().(proto.Message)

	// 將anypb.Any解封裝到創建的實例中
	if err := anyMessage.UnmarshalTo(instance); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Any to target type: %w", err)
	}

	return instance, nil
}
