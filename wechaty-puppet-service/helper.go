package puppetservice

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func grpcTimestampToGoTime(t *timestamppb.Timestamp) time.Time {
	// 不同的 puppet 返回的时间格式不一致, 需要做转换兼容
	// padlocal 返回的是秒，puppet-service 当作是毫秒单位转为秒(除以1000)，所以这里 t.Seconds 就只剩下七位，另外3为被分配到 t.Nanos 去了
	// https://github.com/wechaty/puppet-service/blob/4de1024ee9b615af6c44674f684a84dd8c11ae9e/src/pure-functions/timestamp.ts#L7-L17

	// 这里我们判断 t.Seconds 是否为7位来特殊处理
	//TODO(dchaofei): 未来时间戳每增加一位这里就要判断加一位，那就是200多年之后的事情了，到时还有人在用 wechaty 吗？(2023-09-09)
	if t.Seconds/10000000 < 1 {
		second := t.Seconds*1000 + int64(t.Nanos)/1000000
		return time.Unix(second, 0)
	}

	return t.AsTime().Local()
}
