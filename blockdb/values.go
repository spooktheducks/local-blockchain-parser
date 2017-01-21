package blockdb

type Satoshis int64

func (s Satoshis) ToBTC() BTC {
	return BTC(float64(s) * 0.00000001)
}

type BTC float64

func (b BTC) ToSatoshis() Satoshis {
	return Satoshis(b * 100000000)
}
