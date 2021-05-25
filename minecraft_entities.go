package wrapper

type minecraftEntityID string

func (id minecraftEntityID) Namespace() string {
	return "minecraft"
}

func (id minecraftEntityID) Name() string {
	return string(id)
}

func (id minecraftEntityID) String() string {
	return id.Namespace() + ":" + id.Name()
}

// Constants of minecraft entity IDs
const (
	AreaEffectCloudEntity      minecraftEntityID = "area_effect_cloud"
	ArmorStandEntity           minecraftEntityID = "armor_stand"
	ArrowEntity                minecraftEntityID = "arrow"
	BatEntity                  minecraftEntityID = "bat"
	BeeEntity                  minecraftEntityID = "bee"
	BlazeEntity                minecraftEntityID = "blaze"
	BoatEntity                 minecraftEntityID = "boat"
	CatEntity                  minecraftEntityID = "cat"
	CaveSpiderEntity           minecraftEntityID = "cave_spider"
	ChestMinecartEntity        minecraftEntityID = "chest_minecart"
	ChickenEntity              minecraftEntityID = "chicken"
	CodEntity                  minecraftEntityID = "cod"
	CommandBlockMinecartEntity minecraftEntityID = "command_block_minecart"
	CowEntity                  minecraftEntityID = "cow"
	CreeperEntity              minecraftEntityID = "creeper"
	DolphinEntity              minecraftEntityID = "dolphin"
	DonkeyEntity               minecraftEntityID = "donkey"
	DragonFireballEntity       minecraftEntityID = "dragon_fireball"
	DrownedEntity              minecraftEntityID = "drowned"
	EggEntity                  minecraftEntityID = "egg"
	ElderGuardianEntity        minecraftEntityID = "elder_guardian"
	EndCrystalEntity           minecraftEntityID = "end_crystal"
	EnderDragonEntity          minecraftEntityID = "ender_dragon"
	EnderPearlEntity           minecraftEntityID = "ender_pearl"
	EndermanEntity             minecraftEntityID = "enderman"
	EndermiteEntity            minecraftEntityID = "endermite"
	EvokerEntity               minecraftEntityID = "evoker"
	EvokerFangsEntity          minecraftEntityID = "evoker_fangs"
	ExperienceBottleEntity     minecraftEntityID = "experience_bottle"
	ExperienceOrbEntity        minecraftEntityID = "experience_orb"
	EyeOfEnderEntity           minecraftEntityID = "eye_of_ender"
	FallingBlockEntity         minecraftEntityID = "falling_block"
	FireballEntity             minecraftEntityID = "fireball"
	FireworkRocketEntity       minecraftEntityID = "firework_rocket"
	FishingBobberEntity        minecraftEntityID = "fishing_bobber"
	FoxEntity                  minecraftEntityID = "fox"
	FurnaceMinecartEntity      minecraftEntityID = "furnace_minecart"
	GhastEntity                minecraftEntityID = "ghast"
	GiantEntity                minecraftEntityID = "giant"
	GuardianEntity             minecraftEntityID = "guardian"
	HopperMinecartEntity       minecraftEntityID = "hopper_minecart"
	HorseEntity                minecraftEntityID = "horse"
	HuskEntity                 minecraftEntityID = "husk"
	IllusionerEntity           minecraftEntityID = "illusioner"
	IronGolemEntity            minecraftEntityID = "iron_golem"
	ItemEntity                 minecraftEntityID = "item"
	ItemFrameEntity            minecraftEntityID = "item_frame"
	LeashKnotEntity            minecraftEntityID = "leash_knot"
	LightningBoltEntity        minecraftEntityID = "lightning_bolt"
	LlamaEntity                minecraftEntityID = "llama"
	LlamaSpitEntity            minecraftEntityID = "llama_spit"
	MagmaCubeEntity            minecraftEntityID = "magma_cube"
	MinecartEntity             minecraftEntityID = "minecart"
	MooshroomEntity            minecraftEntityID = "mooshroom"
	MuleEntity                 minecraftEntityID = "mule"
	OcelotEntity               minecraftEntityID = "ocelot"
	PaintingEntity             minecraftEntityID = "painting"
	PandaEntity                minecraftEntityID = "panda"
	ParrotEntity               minecraftEntityID = "parrot"
	PhantomEntity              minecraftEntityID = "phantom"
	PigEntity                  minecraftEntityID = "pig"
	PillagerEntity             minecraftEntityID = "pillager"
	PlayerEntity               minecraftEntityID = "player"
	PolarBearEntity            minecraftEntityID = "polar_bear"
	PotionEntity               minecraftEntityID = "potion"
	PufferfishEntity           minecraftEntityID = "pufferfish"
	RabbitEntity               minecraftEntityID = "rabbit"
	RavagerEntity              minecraftEntityID = "ravager"
	SalmonEntity               minecraftEntityID = "salmon"
	SheepEntity                minecraftEntityID = "sheep"
	ShulkerEntity              minecraftEntityID = "shulker"
	ShulkerBulletEntity        minecraftEntityID = "shukler_bullet"
	SilverfishEntity           minecraftEntityID = "silverfish"
	SkeletonEntity             minecraftEntityID = "skeleton"
	SkeletonHorseEntity        minecraftEntityID = "skeleton_horse"
	SlimeEntity                minecraftEntityID = "slime"
	SmallFireballEntity        minecraftEntityID = "small_fireball"
	SnowGolemEntity            minecraftEntityID = "snow_golem"
	SnowballEntity             minecraftEntityID = "snowball"
	SpawnerMinecartEntity      minecraftEntityID = "spawner_minecart"
	SpectralArrowEntity        minecraftEntityID = "spectral_arrow"
	SpiderEntity               minecraftEntityID = "spider"
	SquidEntity                minecraftEntityID = "squid"
	TNTEntity                  minecraftEntityID = "tnt"
	TNTMinecartEntity          minecraftEntityID = "tnt_minecart"
	TraderLlamaEntity          minecraftEntityID = "trader_llama"
	TridentEntity              minecraftEntityID = "trident"
	TropicalFishEntity         minecraftEntityID = "tropical_fish"
	TurtleEntity               minecraftEntityID = "turtle"
	VexEntity                  minecraftEntityID = "vex"
	VillagerEntity             minecraftEntityID = "villager"
	VindicatorEntity           minecraftEntityID = "vindicator"
	WanderingTraderEntity      minecraftEntityID = "wandering_trader"
	WitchEntity                minecraftEntityID = "witch"
	WitherEntity               minecraftEntityID = "wither"
	WitherSkeletonEntity       minecraftEntityID = "wither_skeleton"
	WitherSkullEntity          minecraftEntityID = "wither_skull"
	WolfEntity                 minecraftEntityID = "wolf"
	ZombieEntity               minecraftEntityID = "zombie"
	ZombieHorseEntity          minecraftEntityID = "zombie_horse"
	ZombiePigmanEntity         minecraftEntityID = "zombie_pigman"
	ZombieVillagerEntity       minecraftEntityID = "zombie_villager"
)

type minecraftEntityTypeTag string

func (tag minecraftEntityTypeTag) Namespace() string {
	return "minecaft"
}

func (tag minecraftEntityTypeTag) Name() string {
	return string(tag)
}

func (tag minecraftEntityTypeTag) String() string {
	return "#" + tag.Namespace() + ":" + tag.Name()
}

// Constants of minecraft entity tags
const (
	Arrows            minecraftEntityTypeTag = "arrows"
	BeehiveInhabitors minecraftEntityTypeTag = "beehive_inhabitors"
	ImpactProjectiles minecraftEntityTypeTag = "impact_projectiles"
	Raiders           minecraftEntityTypeTag = "raiders"
	Skeletons         minecraftEntityTypeTag = "skeletons"
)
