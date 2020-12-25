package wrapper

import (
	"math"
	"time"
)

const (
	// MarketOpenTick is the mc time at which a villagers begin their workday,
	// hence are open for item trading.
	MarketOpenTick int64 = 2000
	// MarketCloseTick is the mc time at which villagers end their workday
	// and begin socializing, trading is not available at this point.
	MarketCloseTick int64 = 9000
	// GameTickPerSecond is the minecraft game server tick runs at a fixed
	// rate of 20 ticks per second.
	GameTickPerSecond int = 20

	ClockSyncInterval = 5 * time.Second
)

type Clock struct {
	ticker     *time.Ticker
	syncTicker *time.Ticker
	LastSync   time.Time
	Tick       int
}

func NewClock() *Clock {
	c := &Clock{
		ticker:     time.NewTicker(1 * time.Second),
		syncTicker: time.NewTicker(ClockSyncInterval),
	}
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.Tick += 20
			}
		}
	}()
	return c
}

func (c *Clock) requestSync() <-chan time.Time {
	return c.syncTicker.C
}

func (c *Clock) stop() {
	c.ticker.Stop()
}

func (c *Clock) resetLastSync() {
	c.LastSync = time.Now()
}

func (c *Clock) syncTick(t int) {
	delay := time.Since(c.LastSync).Seconds()
	delayRoundUp := int(math.Floor(delay))
	tickOffset := delayRoundUp * GameTickPerSecond
	c.Tick = t + tickOffset
}
