# monster-siren-downloader

下载塞壬唱片官网提供的原始音频文件。

Download raw files from Monster Siren Records.

> 塞壬唱片MSR (Monster Siren Records)，泰拉世界十一世纪最大的音乐发行商之一。
>
> 从异教重金属乐队至偶像产业，塞壬唱片MSR旗下的签约艺人涉猎音乐各方面。
>
> 根据最新统计数据，塞壬唱片MSR占有泰拉世界30%以上的音乐市场。
>
> Monster Siren Records (MSR), one of the world's largest music publishers in the eleventh century in Terra.
>
> From heavy metal bands to the idol industry, MSR artists are involved in all aspects of music.
>
> According to the latest statistics, MSR occupies more than 30% of the music market in Terra.

## 使用 / Usage

前往 [Release](https://github.com/AyakuraYuki/monster-siren-downloader/releases) 下载对应操作系统的可执行程序包，解压后运行可执行程序：

- Windows 系统双击 exe 文件即可
- Linux / Unix / macOS 系统需要在可执行程序的当前目录打开终端，运行命令 `./msr-downloader`

音乐会下载保存在当前运行目录下的 `monster-siren` 文件夹内。

---

Go to [Release](https://github.com/AyakuraYuki/monster-siren-downloader/releases) to download the executable package for your operating system. After extracting the files, run the executable:

- Windows: Double-click the `.exe` file
- Linux/Unix/macOS: Open a terminal in the executable's directory and run: `./msr-downloader`

The music will be saved in a `monster-siren` folder within your current working directory.

### macOS 用户 / To macOS users

这是一个开源项目，我没有从其中获得任何收益，所以我没有付费成为 Apple Developer，并对编译好的 macOS 可执行程序进行公证和签名。

我可以保证这个程序没有做任何多余的事，欢迎对代码进行审查。

您可以在第一次运行程序后，前往 [系统设置 - 隐私与安全性 - 安全性] 下面手动允许程序运行。

或者，您也可以使用下面的命令允许程序运行：

```shell
# 前往 msr-downloader 所在的目录
sudo xattr -r -d com.apple.quarantine msr-downloader
```

---

This is an open-source project, and I have not received any financial benefit from it.
Therefore, I have not paid to become an Apple Developer and have not notarized and signed the compiled macOS executable.

I can assure you that the program does nothing unnecessary, and you are welcome to review the code.

After the first run, you can manually allow the program to run by going to [System Settings - Privacy & Security - Security].

Alternatively, you can use the following command to allow the program to run:

```shell
# Navigate to the directory containing msr-downloader
sudo xattr -r -d com.apple.quarantine msr-downloader
```
