package wrapper

import (
	"context"
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
	// ClockSyncInterval is the interval where the wrapper clock will sync with the
	// game tick rate. The wrapper and game tick can be skewed when the game lags.
	ClockSyncInterval time.Duration = 30 * time.Second
)

// clock represents an internal wrapper clock, meant to be always in sync
// with the running game server clock (sync clock.Tick and game server tick).
type clock struct {
	ticker     *time.Ticker
	syncTicker *time.Ticker
	LastSync   time.Time
	Tick       int
}

func newClock() *clock {
	return &clock{
		ticker:     time.NewTicker(1 * time.Second),
		syncTicker: time.NewTicker(ClockSyncInterval),
	}
}

func (c *clock) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.ticker.C:
				c.Tick += GameTickPerSecond
			}
		}
	}()
}

func (c *clock) stop() {
	c.ticker.Stop()
}

func (c *clock) requestSync() <-chan time.Time {
	return c.syncTicker.C
}

func (c *clock) resetLastSync() {
	c.LastSync = time.Now()
}

func (c *clock) syncTick(t int) {
	delay := time.Since(c.LastSync).Seconds()
	delayRoundUp := int(math.Floor(delay))
	tickOffset := delayRoundUp * GameTickPerSecond
	c.Tick = t + tickOffset
}
