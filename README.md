# geodata

定时将 [Loyalsoldier/v2ray-rules-dat](https://github.com/Loyalsoldier/v2ray-rules-dat) 的 GeoSite、GeoIP DAT 转换为 sing-box SRS。

- 每天 UTC 02:23（北京时间/新加坡时间 10:23）检查最新发布；GitHub 定时任务可能延迟。
- 支持 Actions 页面手动运行 **Update SRS rules**。
- `main` 保存工作流与脚本，`srs` 保存完整分类的转换产物。
- 下载同一个上游 release 的 DAT 和 SHA256 文件，校验后转换；任一环节失败不发布。
- GeoSite 保留域名、后缀、关键词、正则及 `@属性` 分类。没有臆造 `@!属性` 分类。
- 全部产物用 sing-box 的 SRS 读取器检查，并核对源数据分类覆盖。输出为 SRS v1（转换器使用 sing-box 1.9.7），不要求最新 SRS 格式。
- 生成可追溯的 `manifest.json` 和 `SHA256SUMS`；上游版本及产物未变时不产生重复提交。

## 下载

```text
https://raw.githubusercontent.com/eloxt/geodata/srs/geosite/cn.srs
https://raw.githubusercontent.com/eloxt/geodata/srs/geosite/geolocation-!cn.srs
https://raw.githubusercontent.com/eloxt/geodata/srs/geoip/cn.srs
```

分类索引与 SHA256：[manifest.json](https://github.com/eloxt/geodata/blob/srs/manifest.json)。每个分类一个文件，目录分别是 `geosite` 和 `geoip`。只有上游 DAT 中实际存在的类别才会发布。

## sing-box 配置片段

合并到现有 `route.rule_set`；此片段只加载规则集，实际出口由你的路由规则指定。

```json
{
  "tag": "geosite-cn",
  "type": "remote",
  "format": "binary",
  "url": "https://raw.githubusercontent.com/eloxt/geodata/srs/geosite/cn.srs",
  "update_interval": "24h"
}
```

## 实现与验证范围

使用 [runetfreedom/geodat2srs](https://github.com/runetfreedom/geodat2srs) 的固定提交 `368e4b324bb4b8c90da93d1b8429bd0af337dfda`，不会在每次运行时自动升级转换器。工作流使用仓库自带的 `GITHUB_TOKEN`，无需配置个人 Token。

SRS 可读与分类覆盖检查不等同于全部匹配语义的形式化证明。更换转换器或 sing-box 版本时，应抽查关键域名、正则和属性分类。上游若新增当前转换器不支持的 inverse GeoIP，校验将失败，避免发布错误规则。

GitHub 可能在公开仓库长期无活动时停用定时任务；若发现停更，请检查 Actions 是否被禁用。分支保护也不能阻止工作流向 `srs` 分支写入。

上游数据与第三方工具遵循各自许可证；本项目不改变它们的授权。
