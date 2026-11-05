---
name: ctf-forensics
description: CTF Misc/Forensics 方向自主解题：文件类型识别与修复、隐藏数据提取、图片与音频隐写、压缩包伪加密与破解、内存与磁盘取证、流量分析（PCAP）、编码与零宽字符隐写。适用于取证题与 Misc 题。
license: MIT
compatibility: 仅限 CTF / 靶场 / 授权环境
---

# CTF Misc / Forensics 自主解题

## 0. ⚠️ 你看不见图片 —— 必须先把它变成文字

**你没有视觉能力**。用户上传的图片（png/jpg/gif/bmp/webp/tiff）你用 Read 只会读到二进制垃圾，
直接说"我看不了图片"是**错的**。正确做法：用 `imgread` 把图片翻译成文字再看。

```bash
imgread <图片路径>                 # 全量分析（推荐第一步）
imgread <图片路径> --no-ascii      # 只要文字/元数据，不要轮廓图
imgread <图片路径> --quiet         # 只输出 OCR 文本（便于管道处理）
imgread <图片路径> --width 120     # 调大 ASCII 预览宽度，看更细的轮廓
```

`imgread` 一次给出 7 类信息：

| # | 内容 | 对你有什么用 |
|---|---|---|
| 1 | 格式/尺寸/色彩模式/帧数 | 判断是不是多帧图、索引色（调色板题） |
| 2 | EXIF / 元数据 | 相机型号、时间、**注释字段常直接藏 flag** |
| 3 | **OCR 文字**（tesseract 中英文） | 截图题、文字题的关键，直接读出台词/提示 |
| 4 | **ASCII 亮度预览** | **让你"看见"布局**：轮廓、位置、明暗分布 |
| 5 | 颜色分析 | 主色占比、是否高度单色（纯色藏数据/调色板异常） |
| 6 | 隐写快检 | PNG→zsteg/pngcheck；尾部附加数据；LSB 置1比例；strings 命中 |
| 7 | 下一步建议 | 按图片类型给具体命令 |

**判读要点**：
- OCR 出文字 → 先信 OCR 结果，往往直接就是 flag 或提示
- OCR 空 → 放大再试：`convert in.png -resize 400% -colorspace Gray /tmp/big.png && imgread /tmp/big.png`
- ASCII 预览里看到规整方块/条纹/QR 样式 → 大概率是二维码或图案编码，走 `zbarimg`
- 颜色分析报"高度单色" → 优先查调色板/LSB/纯色区域藏数据
- EXIF 有注释/经纬度/软件名 → 直接看，CTF 老套路

**现在有 Pillow / opencv / scikit-image / scipy / pytesseract 了**，
不要再把图转成 .pgm/.npy 硬啃（除非题目本身就要求像素级分析）：

```python
from PIL import Image
im = Image.open('x.png'); print(im.size, im.mode)
# 多帧
for i in range(getattr(im, 'n_frames', 1)): im.seek(i); im.convert('RGB').save(f'f{i}.png')
# 调色板
print(im.getpalette()[:30])
```

## 1. 拿到文件先做（顺序不可颠倒）

```bash
cd /app/data/tasks/TID/ && mkdir -p forensics && cd forensics
f=./attachment
file "$f"; stat -c '%s bytes' "$f"      # 1) 真实类型（别信扩展名）2) 大小
xxd "$f" | head -20                      # 3) 文件头 / 尾部魔数
binwalk "$f"                             # 4) 嵌套 / 追加数据
exiftool "$f" | head -40                 # 5) 元数据（注释 / 作者 / 经纬度）
strings -a -n 6 "$f" | grep -iE 'flag|\{' | head   # 6) 明文字符串
# 图片再补一步：
case "$(file -b --mime-type "$f")" in image/*) imgread "$f" ;; esac
```

**扩展名与 `file` 结果不一致 → 先改名按真实类型处理**（最常见的入口）。

## 2. 分类型处理

| 类型 | 关键动作 |
|---|---|
| **PNG** | `pngcheck -v` 查完整性；IHDR 宽高被改 → 修复；**注意非标准块**（如 `hIDn`：里面常是 [4B宽][4B高]+zlib 的被裁部分，`imgread` 会自动提取）；`binwalk -e` 提追加数据；`zsteg` 扫 LSB |
| **JPEG** | 尾部追加数据、EXIF 注释、缩略图差异 |
| **GIF** | 逐帧拆（`convert -coalesce`），帧间差异藏数据 |
| **ZIP** | 伪加密（改回 flag 位）、注释、CRC 爆破、`zip2john` + `john` 破密码 |
| **PCAP** | `tshark -r x.pcap -Y http` 提对象、`-z follow,tcp,ascii,N` 跟流、导出文件 |
| **内存镜像** | `volatility3 -f mem.raw windows.info`，再按插件找进程 / 文件 / 剪贴板 |
| **磁盘镜像** | `mmls` / `fls` / `icat`（sleuthkit）或 `testdisk` 恢复分区与文件 |
| **音频** | 看频谱（`sox`）；摩斯 / DTMF 解码；WAV 最低位 |
| **纯文本** | 零宽字符、不可见字符、多重编码、`\r\n` 差异 |
| **二维码** | `zbarimg` 扫；可修复破损二维码 |

```bash
zsteg -a "$f" 2>/dev/null | head -40
steghide extract -sf "$f" -p '' 2>/dev/null
stegseek "$f" ${WORDLISTS:-/opt/ctf/wordlists}/rockyou.txt 2>/dev/null
foremost -i "$f" -o out_foremost 2>/dev/null
zip2john x.zip > h.txt && john --wordlist=${WORDLISTS:-/opt/ctf/wordlists}/rockyou.txt h.txt
tshark -r x.pcap --export-objects http,out_http 2>/dev/null && ls out_http
```

## 3. 两文件对比（高频技巧）

```bash
cmp -l a.png b.png | head
python3 - <<'PY'
a = open('a.png', 'rb').read(); b = open('b.png', 'rb').read()
diff = [(i, x, y) for i, (x, y) in enumerate(zip(a, b)) if x != y]
print(len(diff), diff[:20])
print(''.join(chr(y) for _, _, y in diff if 32 <= y < 127))   # 差异字节里常藏明文
PY
```

## 4. 纪律

- **没看到文件就不要说"这里有 LSB 隐写"**；先 `file` / `xxd` / `binwalk` 拿到证据。
- 每层中间产物都落盘，并检查它是否又是另一种格式（递归处理）。
- 候选 flag 按 `ctf-solve` 的判真协议验证（题目可能埋 decoy）。
