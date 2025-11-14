/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package SetDictOneStop
 *@file    init
 *@date    2025/7/15 21:33
 */

package SetDictOneStop

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var itemCache = cache.New(24*time.Hour, 24*time.Hour)

func init() {

}
