package broker

import (
	"reflect"
	"testing"
)

func TestEventQueue_HandleEvent(t *testing.T) {
	type args struct {
		eventName string
		opts      []DialOption
	}
	testCases := []struct {
		name string
		args args
		want *eventQueue
	}{
		{
			name: "default options for queue name and route key",
			args: args{
				eventName: "user.created",
			},
			want: &eventQueue{
				events: map[string]*handler{
					"user.created": &handler{
						Options: dialOptions{
							QueueName: "user.created.queue",
							RouteKey:  "",
						},
					},
				},
			},
		},
		{
			name: "with queue name and route key defualt",
			args: args{
				eventName: "user.created",
				opts:      []DialOption{WithQueue("user.created.my-queue")},
			},
			want: &eventQueue{
				events: map[string]*handler{
					"user.created": &handler{
						Options: dialOptions{
							QueueName: "user.created.my-queue",
							RouteKey:  "",
						},
					},
				},
			},
		},
		{
			name: "with queue name and route key defualt",
			args: args{
				eventName: "user.created",
				opts:      []DialOption{WithRouteKey("new_user"), WithQueue("user.created.my-queue")},
			},
			want: &eventQueue{
				events: map[string]*handler{
					"user.created": &handler{
						Options: dialOptions{
							QueueName: "user.created.my-queue",
							RouteKey:  "new_user",
						},
					},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			eq := &eventQueue{
				events: map[string]*handler{},
			}
			eq.HandleEvent(tc.args.eventName, nil, tc.args.opts...)
			if !reflect.DeepEqual(eq, tc.want) {
				t.Errorf("HandleEvent(event,handler,opts) got %v; want %v", eq, tc.want)
			}
		})
	}
}
