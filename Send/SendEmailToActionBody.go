/*
 *@author  chengkenli
 *@project StarRocksDict
 *@package Send
 *@file    SendEmailToActionBody
 *@date    2024/8/30 16:09
 */

package Send

import (
	"StarRocksDict/util"
	"fmt"
)

func HtmlEmailBody(dict *util.DictOption) ([]string, string) {
	var user []string
	if dict.User != "" {
		user = append(user, dict.User+"@xxx.com")
	}
	if dict.Owner != "" {
		user = append(user, dict.Owner+"@xxx.com")
	}

	var action string
	switch dict.Option {
	case "create":
		action = "<strong>创建</strong>"
	case "drop":
		action = "<strong>删除</strong>"
	case "reset":
		action = "<strong>重置密码</strong>"
	}

	zhuyi := fmt.Sprintf(`
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	请注意，登录时<strong>账号严格区分大小写</strong>，请您<strong>使用小写进行登录</strong>。
</p>
`)
	var username, password string
	username = dict.User
	if dict.Option == "drop" {
		password = "已销毁"
		zhuyi = ""
	} else {
		password = dict.Password
	}

	msg := fmt.Sprintf(`
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	主题：StarRocks数据库账号%s成功！
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<br />
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	尊敬的用户，
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	您好！您的StarRocks数据库账号已%s。
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<br />
</p>
<ul>
	<li>
		【账号】：[<strong><span style="color:#009900;">%s</span></strong>]
	</li>
	<li>
		【密码】：[<strong>%s</strong>]
	</li>
	<li>
		【集群】：<strong>%s</strong> 
	</li>
	<li>
		【网页】：%s
	</li>
	<li>
		【其他】：%s
	</li>
</ul>
<p>
	<br />
</p>
%s
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	为了帮助您更好地使用StarRocks数据库，我们为您提供了以下参考资料：
</p>
<ol style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<li>
		<p>
			登录指引：请参考以下链接，了解如何登录StarRocks数据库并进行基本操作。 [<a href="xxx" target="_blank">登录指引链接</a>]
		</p>
	</li>
	<li>
		<p>
			飞书群组：了解更多关于StarRocks的功能、日常支持、流程指引、答疑解惑。 [<a href="xxx" target="_blank">加入飞书群组</a>]
		</p>
	</li>
	<li>
		<p>
			技术文档：在飞书群<strong>群公告</strong>中存在大量技术指引可满足您的需求。
		</p>
	</li>
</ol>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<br />
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	如有其他疑问或需要帮助，请随时联系我们。感谢您使用StarRocks！
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	此电子邮件由系统自动发送，请勿直接回复此电子邮件。&nbsp;
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<br />
</p>
<p style="font-family:&quot;font-size:15px;background-color:#FFFFFF;">
	<strong>Infra Data Team</strong> 
</p>`, action, action, username, password, dict.StarRocks, dict.Web, dict.Other, zhuyi)

	return user, msg
}
