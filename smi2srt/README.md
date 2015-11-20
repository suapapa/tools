# smi2srt: Script format changer

Install:

    go get github.com/suapapa/tools/smi2srt

Example Usage:

    smi2srt -o output.srt input_utf8.smi

The `smi` file should be encoded in `EUC-KR` or `UTF-8`.

`smi2srt` detect input encoding with `github.com/saintfish/chardet`
and, if its `EUC-KR` decode it to `UTF-8` with
`github.com/suapapa/go_hangul/encoding/cp949`.

## livecoding videos
* [1/4](https://youtu.be/lY_SQC_oFPM)
* [2/4](https://youtu.be/GuYWOmbsK-4)
* [3/4](https://youtu.be/xMxUEoYWcQ8)
* [4/4](https://youtu.be/s3my9TX7wUE)
