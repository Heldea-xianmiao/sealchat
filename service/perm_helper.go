package service

import (
	"strings"

	"sealchat/model"
)

// ResolveMemberRoleForProtocol returns the protocol role for a user in a channel/world context.
// Returns: "owner" | "admin" | "member"
func ResolveMemberRoleForProtocol(userID, channelID, worldID string) string {
	if strings.TrimSpace(userID) == "" {
		return model.WorldRoleMember
	}

	// 1. Check channel ownership and channel roles
	if strings.TrimSpace(channelID) != "" {
		channel, err := model.ChannelGet(channelID)
		if err == nil && channel != nil {
			if channel.UserID != "" && channel.UserID == userID {
				return model.WorldRoleOwner
			}
			if strings.TrimSpace(worldID) == "" && strings.TrimSpace(channel.WorldID) != "" {
				worldID = channel.WorldID
			}
		}

		// 1.5 Check channel role (ch-{channelID}-owner or ch-{channelID}-admin)
		roleIDs, err := model.UserRoleMappingListByUserID(userID, channelID, "channel")
		if err == nil && len(roleIDs) > 0 {
			for _, roleID := range roleIDs {
				if strings.HasSuffix(roleID, "-owner") {
					return model.WorldRoleOwner
				}
			}
			for _, roleID := range roleIDs {
				if strings.HasSuffix(roleID, "-admin") {
					return model.WorldRoleAdmin
				}
			}
		}
	}

	if strings.TrimSpace(worldID) == "" {
		return model.WorldRoleMember
	}

	// 2. Check world ownership
	world, err := GetWorldByID(worldID)
	if err == nil && world != nil && world.OwnerID != "" && world.OwnerID == userID {
		return model.WorldRoleOwner
	}

	// 3. Check world member role
	var member model.WorldMemberModel
	if err := model.GetDB().Where("world_id = ? AND user_id = ?", worldID, userID).Limit(1).Find(&member).Error; err == nil && member.ID != "" {
		switch member.Role {
		case model.WorldRoleOwner:
			return model.WorldRoleOwner
		case model.WorldRoleAdmin:
			return model.WorldRoleAdmin
		}
	}

	return model.WorldRoleMember
}
