## v1.1.0

- 更新到 Go 1.26
- 使用 [monster-siren-api-go](https://github.com/AyakuraYuki/monster-siren-api-go) 包装的API SDK获取专辑和歌曲的数据
- 优化了组装专辑信息文本的代码
- 修复了显示进度条时，进度条伸缩的问题
- 支持显示当前正在下载的文件的下载进度

---

- Update to Go 1.26.
- Use the [monster-siren-api-go](https://github.com/AyakuraYuki/monster-siren-api-go) API SDK to fetch album and song data.
- Optimized the code for assembling album information text.
- Fixed an issue where the progress bar would stretch when displayed.
- Added support for displaying the download progress of the current file being downloaded.

## v1.0.9

- 运行优化
- 恢复为不覆盖`info.txt`，因为有反病毒软件以此将上一个版本误判成有病毒的程序

---

- Performance optimisations.
- Reverted to not overwriting `info.txt` because some antivirus software falsely detected the previous version as a virus.

## v1.0.8

- 修复了保存专辑信息时丢失 `专辑作者` 的问题
- 总是保存最新的专辑信息到 `info.txt` 文件

---

- Fixed an issue where the `专辑作者` was lost when saving album information.
- Always save the latest album information to the `info.txt` file.

## v1.0.7

- 修复了在获取专辑的文件夹名称前，没有修剪前后空白字符的问题
- 因Windows系统对结尾以`.`号命名的文件夹并不友好，导致程序不能识别到路径，将专辑名称结尾字符是`.`号的，替换成以`_`结尾

---

- Fixed an issue where the folder name of an album was not trimmed of leading and trailing whitespace characters before being retrieved.
- Due to Windows systems being unfriendly to folder names ending with a period (`.`), which prevents the program from recognizing the path, album names ending with a period are now replaced with an underscore (`_`) at the end.

## v1.0.6 in Pre-release

- 修复了 Windows 系统不能使用 `.` 号作为文件夹名称结尾的问题

---

- Fixed unsupported character in filename suffix in Windows.

## v1.0.5

- 保存临时文件以支持进程中断后重新下载

---

- Support for resuming downloads after unexpected program termination and restart by using temporary file.

## v1.0.4 in Pre-release

- 修剪了专辑文件夹的前后空白字符，部分专辑因为名称前后存在空白字符，增量下载不能识别，因此会重新下载这些专辑
- 修复了进度条在下载结束后不能正常输出的问题

---

- Trimmed the leading and trailing blank characters of the album folder. Some albums have blank characters before or after the name, so the incremental download cannot recognize them and will
  re-download these albums.
- Fixed the problem that the progress bar could not be output normally after the download was completed.

## v1.0.3 in Pre-release

> Version 1.0.2 disappeared for some embarrassing reasons.

- Support download progress.
- Show the logs in progress.

## v1.0.1

- 支持最大同时5个下载任务（不可调整）
- 延长了请求超时，规避因为超时导致下载音乐失败

---

- Supports a maximum of 5 concurrent download tasks (non-adjustable)
- Extended request timeout to prevent music download failures due to timeouts

## v1.0.0

- The first version of this downloader
- Downloads songs one by one (not in multitasks)
