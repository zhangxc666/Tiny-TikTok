## Tiny Tiktok

主要根据[简单代码](https://github.com/My-younth-is-over/Easy-version-Tiktok)结构基础上，学习[第六届青训营第一名代码](https://github.com/Happy-Why/toktik)代码进行修改优化。

### 技术栈

Gin+Gorm+Go-Redis+Ffmpeg+Mysql+Redis

### 完成列表

- 完成message功能，增加了message缓存（写回）
- 完成User功能，增加了count（写回），info（旁路缓存）的缓存
- 完成follow功能，增加了关注、粉丝和好友id的缓存（延迟双删）
- 完成video功能，增加了video_info（写回）缓存和系统用户发布视频（旁路缓存）缓存
  - [ ] 上传至OSS中
  - [ ] 增加video计数表解耦  
- 完成favor功能，增加了点赞（写回）缓存

