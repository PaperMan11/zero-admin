package utils

import (
	"reflect"
	"testing"
)

func TestMapToStruct(t *testing.T) {
	type ChatRoomLog struct {
		Uid     int                    `json:"uid"`
		Online  int                    `json:"online"`   // -1 注销
		RoomId  int                    `json:"room_id"`  // 房间id
		OwnerId int                    `json:"owner_id"` // 房主id
		Action  string                 `json:"action"`
		Time    int                    `json:"time"`
		Ext     map[string]interface{} `json:"ext"`
	}

	tests := []struct {
		name       string
		input      map[string]interface{}
		want       ChatRoomLog
		shouldFail bool
	}{
		{
			name: "AllFieldsProvided_ShouldMapSuccessfully",
			input: map[string]interface{}{
				"uid":      123,
				"online":   1,
				"room_id":  456,
				"owner_id": 789,
				"action":   "join",
				"time":     1620000000,
				"ext": map[string]interface{}{
					"device": "mobile",
					"ip":     "192.168.1.1",
				},
			},
			want: ChatRoomLog{
				Uid:     123,
				Online:  1,
				RoomId:  456,
				OwnerId: 789,
				Action:  "join",
				Time:    1620000000,
				Ext: map[string]interface{}{
					"device": "mobile",
					"ip":     "192.168.1.1",
				},
			},
			shouldFail: false,
		},
		{
			name: "SomeFieldsMissing_ShouldMapRemaining",
			input: map[string]interface{}{
				"uid":    123,
				"action": "join",
			},
			want: ChatRoomLog{
				Uid:    123,
				Action: "join",
			},
			shouldFail: false,
		},
		{
			name: "FieldTypeMismatch_ShouldReturnError",
			input: map[string]interface{}{
				"uid": "not_an_int", // 类型错误
			},
			want:       ChatRoomLog{},
			shouldFail: true,
		},
		{
			name: "ExtraFieldIgnored_ShouldMapSuccessfully",
			input: map[string]interface{}{
				"uid":      123,
				"extraKey": "value",
			},
			want: ChatRoomLog{
				Uid: 123,
			},
			shouldFail: false,
		},
		{
			name: "ExtIsMapWithNestedValues_ShouldMapSuccessfully",
			input: map[string]interface{}{
				"uid": 123,
				"ext": map[string]interface{}{
					"key1": "val1",
					"key2": map[string]interface{}{"subkey": "subval"},
				},
			},
			want: ChatRoomLog{
				Uid: 123,
				Ext: map[string]interface{}{
					"key1": "val1",
					"key2": map[string]interface{}{"subkey": "subval"},
				},
			},
			shouldFail: false,
		},
		{
			name: "ExtIsNotAMap_ShouldFail",
			input: map[string]interface{}{
				"uid": 123,
				"ext": "this_is_not_a_map",
			},
			want:       ChatRoomLog{},
			shouldFail: true,
		},
		{
			name: "PointerTypeFieldWithValue_ShouldMapSuccessfully",
			input: map[string]interface{}{
				"ext": map[string]interface{}{
					"test": "value",
				},
			},
			want: ChatRoomLog{
				Ext: map[string]interface{}{
					"test": "value",
				},
			},
			shouldFail: false,
		},
		{
			name:  "EmptyInput_ShouldReturnZeroValue",
			input: map[string]interface{}{
				// 空输入
			},
			want:       ChatRoomLog{},
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got ChatRoomLog
			err := MapToStruct(tt.input, &got)

			if (err != nil) != tt.shouldFail {
				t.Fatalf("MapToStruct() error = %v, wantErr %t", err, tt.shouldFail)
			}

			if !tt.shouldFail && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MapToStruct() got = %+v, want = %+v", got, tt.want)
			}
		})
	}
}
