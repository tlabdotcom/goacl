package goacl

import (
	"context"
	"encoding/json"

	"github.com/labstack/gommon/log"
)

type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

func (a *ACL) publishEvent(eventType string, payload interface{}) {
	event := Event{
		Type:    eventType,
		Payload: payload,
	}
	jsonEvent, _ := json.Marshal(event)
	a.Redis.RPush(context.Background(), a.AclKeyEvents, jsonEvent)
}

func (a *ACL) ConsumeEvents(ctx context.Context) {
	for {
		// Blocking call to retrieve the event from Redis
		events, err := a.Redis.BLPop(ctx, 0, a.AclKeyEvents).Result()
		if err != nil {
			log.Errorf("Error consuming events: %v", err)
			continue
		}

		// The event data is in the second element of the returned slice
		if len(events) < 2 {
			continue
		}

		var event Event
		if err := json.Unmarshal([]byte(events[1]), &event); err != nil {
			log.Errorf("Error unmarshalling event: %v", err)
			continue
		}

		switch event.Type {
		case RoleCreatedEvent:
			role := event.Payload.(*AclParam)
			err = a.createRoleToDB(ctx, role)
			if err != nil {
				log.Errorf("Error adding role: %v", err)
			}
		case RoleDeletedEvent:
			role := event.Payload.(*Role)
			err = a.updateRoleToDB(ctx, role)
			if err != nil {
				log.Errorf("Error updating role: %v", err)
			}

		case RoleUpdatedEvent:
			role := event.Payload.(*Role)
			err = a.deleteRoleToDB(ctx, role)
			if err != nil {
				log.Errorf("Error deleting role: %v", err)
			}

		case FeatureCreatedEvent:
			feature := event.Payload.(*Feature)
			err = a.createFeatureToDB(ctx, feature)
			if err != nil {
				log.Errorf("Error inserting feature: %v", err)
			}

		case FeatureUpdatedEvent:
			feature := event.Payload.(*Feature)
			err := a.updateFeatureToDB(ctx, feature)
			if err != nil {
				log.Errorf("Error updating feature: %v", err)
			}

		case FeatureDeletedEvent:
			feature := event.Payload.(*Feature)
			err := a.deleteFeatureToDB(ctx, feature)
			if err != nil {
				log.Errorf("Error deleting feature: %v", err)
			}

		case SubFeatureCreatedEvent:
			subFeature := event.Payload.(*SubFeature)
			err = a.createSubFeatureToDB(ctx, subFeature)
			if err != nil {
				log.Errorf("Error inserting subfeature: %v", err)
			}

		case SubFeatureUpdatedEvent:
			subFeature := event.Payload.(*SubFeature)
			err := a.updateSubFeatureToDB(ctx, subFeature)
			if err != nil {
				log.Errorf("Error updating subfeature: %v", err)
			}

		case SubFeatureDeletedEvent:
			subFeature := event.Payload.(*SubFeature)
			err := a.deleteSubFeatureToDB(ctx, subFeature)
			if err != nil {
				log.Errorf("Error deleting subfeature: %v", err)
			}

		case PolicyCreatedEvent:
			policy := event.Payload.(*Policy)
			err := a.createPolicyToDB(ctx, policy)
			if err != nil {
				log.Errorf("Error inserting policy: %v", err)
			}

		case PolicyDeletedEvent:
			policy := event.Payload.(*Policy)
			err := a.deletePolicyToDB(ctx, policy)
			if err != nil {
				log.Errorf("Error deleting policy: %v", err)
			}
		}
	}
}
