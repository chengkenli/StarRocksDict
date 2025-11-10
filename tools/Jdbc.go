/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package tools
 *@file    Jdbc
 *@date    2024/8/30 16:43
 */

package tools

import "StarRocksDict/util"

func JdbcHost(app string) (string, string) {
	var web, jdbc string
	for _, m := range util.MetaLink {
		if m["app"].(string) == app {
			web = m["address"].(string)
			jdbc = m["feip"].(string)
			break
		}
	}
	return web, jdbc
}
