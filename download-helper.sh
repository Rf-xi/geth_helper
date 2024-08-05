#!/bin/bash

# 设置 Tmux 目标窗格名称
control_pane="6"
geth_pane="4"
prysm_pane="5"

function get_latest_output() {
    # 查询最新块的日期
    tmux send-keys -t $control_pane 'new Date(eth.getBlock("latest").timestamp * 1000).toLocaleString()' Enter 
    
    # 等待输出
    sleep 1

    # 获取最后两行输出
    tmux capture-pane -p -t $control_pane | tail -n 2
}

function stop_tmux_program() {
    # 发送 Ctrl+C 命令停止 Tmux 程序
    tmux send-keys -t $geth_pane  C-c Enter
    tmux send-keys -t $prysm_pane C-c Enter

    echo "Stopped the Tmux program."
}

# 检查命令行参数
if [ "$1" == "query" ]; then
    # 输出最新的输出信息
    get_latest_output
elif [ "$1" == "stop" ]; then
    # 停止 Tmux 程序
    stop_tmux_program
else
    echo "Usage: $0 [query|stop]"
    exit 1
fi
