# -*- coding: UTF-8 -*-
# windows下的编译脚本
# author: wangfeng
# since: 2023-09-12

import os
import platform
import subprocess
from git import Repo


def get_latest_tag(diff: int=0) -> str:
    """
    获取仓库最新的一个tag, 支持传入调整的数值, 默认只修改修订版本号
    """
    repo = Repo(r'./')
    tags = []
    for __tag in repo.tags:
        tag = str(__tag)
        tag = tag[1:]
        tags.append(tag)
    tags.sort(key=lambda x:tuple(int(v) for v in x.split('.')))
    latest = tags[-1]
    if diff != 0:
        last_vs = tags[-1].split('.')
        last_vs[-1] = str(int(last_vs[-1])+diff)
        latest = '.'.join(last_vs)
    return latest


# 只能获取执行结果
def subprocess_check_output(stmt, shell:bool=False):
    result = subprocess.check_output(stmt, shell).decode('utf-8')
    # 执行失败不需要特殊处理，命令执行失败会直接报错
    return result  # 返回执行结果，但是结果返回的是一个str字符串（不论有多少行），并且返回的结果需要转换编码

def fix_version(v: str) -> str:
    if v.startswith(('v','V')):
        v = v[1:]
    return v

def get_mod_version(module: str) -> str:
    # 获取依赖库stock的版本号
    cmd_result = subprocess_check_output(f'go list -m {module}').strip('\n')
    vs = cmd_result.split(' ')
    tag_latest = fix_version(vs[1])
    version = tag_latest
    return version

if __name__ == '__main__':
    print(os.path.abspath('.'))
    # 获取应用的版本号
    version = get_latest_tag()
    # 获取依赖库的版本号
    gotdx_version = get_mod_version('gitee.com/quant1x/gotdx')
    print('gotdx version: ', gotdx_version)
    repo = 'gitee.com/quant1x/engine'
    BIN='./bin'
    current_path = os.path.dirname(os.path.abspath(__file__))
    #print(current_path)
    BIN=current_path + '/' + BIN
    APP = 'engine'
    EXT = '.exe'
    GOOS = platform.system().lower()
    machine = platform.machine().lower()
    GOARCH = 'amd64' if machine=='amd64' else 'arm64'
    print(f"正在编译应用:{APP} => {BIN}/{APP}{EXT}...")
    cmd= fr'''go env -w GOOS={GOOS} GOARCH={GOARCH} && go build -ldflags "-s -w -X 'main.MinVersion={version}' -X 'main.tdxVersion={gotdx_version}'" -o {BIN}/{APP}{EXT} {repo}'''
    print(cmd)
    subprocess.Popen(cmd, stdout=subprocess.PIPE, shell=True)
    #subprocess.call(cmd,  stdout=subprocess.PIPE, shell=True)
    print(f"正在编译应用:{APP} => {BIN}/{APP}{EXT}...OK")

