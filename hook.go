package zlog

type ErrorHook func(error, *Log)

type LogHook func(*Log)
