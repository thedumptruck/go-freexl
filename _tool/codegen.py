#!/usr/bin/env python2
# -*- coding: utf-8 -*-

import re

opts = []
codes = []
infos = []
pattern = re.compile('^\w+_DECLARE.+')
for line in open("/usr/local/opt/freexl/include/freexl.h"):
    if line.startswith('#define FREEXL_'):
        c = line.split()
        codes.append(c[1])

template = """
// generated by codegen.py

package freexl
/*
#include <freexl.h>
*/
import "C"

// FREEXL
const (
{code_part}
)
// generated ends
"""

code_part = []
for c in codes:
    if c == 'FREEXL_DECLARE':
        continue
    code_part.append("\t{:<26} = C.{}".format(c[7:], c))

code_part = '\n'.join(code_part)

with open('./const_gen.go', 'w') as fp:
    fp.write(template.format(**locals()))
