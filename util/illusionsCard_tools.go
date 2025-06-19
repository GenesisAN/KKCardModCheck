package util

import (
	icb "github.com/GenesisAN/illusionsCard/Base"
)

// 提取卡片缺失mod
func CollectMissingMods(cards []icb.CardInterface, localGUIDs []string) map[string]icb.ResolveInfo {
	missing := make(map[string]icb.ResolveInfo)
	for _, card := range cards {
		if card == nil {
			continue
		}
		if comparer, ok := card.(interface {
			CompareMissingMods([]string) map[string]icb.ResolveInfo
		}); ok {
			for guid, info := range comparer.CompareMissingMods(localGUIDs) {
				missing[guid] = info
			}
		}
	}
	return missing
}
