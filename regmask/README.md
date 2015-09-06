# regmask : print bits fields of given value

[![Build Status](https://travis-ci.org/suapapa/regmask.png?branch=master)](https://travis-ci.org/suapapa/regmask)

# Install

    $ go get github.com/suapapa/regmask

# Usage

Make csv file which describe `bitsFields` in form of `lable,bit_offset,bit_cnt`;

    $ cat example.csv
    IDS_ARM,24,8
    MOD_SG,20,3
    ORIG_SG,17,4
    MPM,12,5
    LOCKING,7,5
    EMA,6,1
    _,4,2
    FUSED_SG,3,1
    MOD_SG,0,3

Use `regmask` to print `bitsFields` of given values;

    $ regmask -g example.csv 0x10094008 0x0804f008
    |      value | IDS_ARM[31:24] | MOD_SG[22:20] | ORIG_SG[20:17] | MPM[16:12] | LOCKING[11:7] | EMA[6] | _[5:4] | FUSED_SG[3] | MOD_SG[2:0] |
    |------------|----------------|---------------|----------------|------------|---------------|--------|--------|-------------|-------------|
    | 0x10094008 |      00010000b |          000b |          0100b |     10100b |        00000b |     0b |    00b |          1b |        000b |
    | 0x0804f008 |      00001000b |          000b |          0010b |     01111b |        00000b |     0b |    00b |          1b |        000b |

Try other base for value with option `-b` and other output format with option '-f'

    $ ./regmask -b 16 -f csv  example.csv 0x10094008 0x0804f008
    value,IDS_ARM,MOD_SG,ORIG_SG,MPM,LOCKING,EMA,_,FUSED_SG,MOD_SG
    0x10094008,0x10,0x0,0x4,0x14,0x00,0x0,0x0,0x1,0x0
    0x0804f008,0x08,0x0,0x2,0x0f,0x00,0x0,0x0,0x1,0x0
