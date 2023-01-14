## Gameplay

The game starts with a game module selection (for instance, a Crypt).

Then players create a party.

The party starts at the entrance of the module. It could be considered to be a first level of a dungeon.

The game loop looks like this:

1. Choose a tile to visit (pick one of 3)
  * Roll an encounter
    * Clear the encounter (if any)
  * Clear the tile event (usually involves dice rolling)
2. If dungeon level is not complete, repeat step 1
3. Fight a level boss
  * Roll a boss fight reward (items, gold)
  * Roll 3 skill cards for every player
4. If it was the last dungeon level, it's victory
5. Return to the city: can rest and buy items now
6. Back to the dungeon: start a new dungeon level, repeat step 1

A dungeon level boss is encountered after clearing the required number of tiles.

Random encounters rarely give any benefits, but scavenge roll is done after every such battle.

The heroes try to survive through the dungeon level and defeat the boss. The dungeon tiles can give them temporary or permanent bonuses. Higher rewards usually mean higher risks. After the boss is defeated, every hero gets a new skill. A party has the chance to resupply at the town and go back to the dungeon, to start clearing the next level.

## Heroes

### Character Creation

1. Roll a starting skill 
2. Choose a starting weapon
3. Choose a talent
4. Choose a starting trait
5. Choose how your character looks like

### Hero Traits

Starting traits:

* Warrior trait: +3 max health
* Mage trait: +5 max energy
* Rogue trait: +2 max health, +3 max energy

### Hero Talents

* Lockpicking (very useful in some tiles)
* Search (bonus roll point in most roll-for-reward tiles)
* Trap disarming (very useful in some tiles)
* Appraise (better buy/sell prices in the city)
* Alchemy (mix potions for greater potions)
* Lore (identify items outside of the city)
* Scavenge (improves the chance of getting any loot after the random encounter)

**Lockpicking** allows a character to open a locked object with 100% chance. Without such talent, there is only ~16.6% of success (6 should be rolled).

**Search** allows to detect a hidden object or room, opening more possibilities for exploration.

**Trap disarming** allows a character to safely disarm any trap. Some random events may have trap-related moments. Some hidden rooms may be guarded by a trapped door. In addition to that, a treasure chest can be protected by a trap in addition to being locked.

**Appraise** improves the prices in the shops. You can sell items for more gold and buy new equipment for lower prices.

**Alchemy** allows potion mixing to achieve unique results. It also improves the potion effects when they're consumed by the character with this talent.

**Lore** allows a party to identify any new item at once, anywhere. Without this talent, items can be identified only in dedicated places like shops.

**Scavenge** adds +2 to a scavenge roll after a random encounter.

## Effects

### Negative Effects

Effects that are simple opposites of positive effects are not listed.

* Skill usage ban (silence)
* Movement action ban (dryad root)
* Attack action ban (disarm)
* Stun
* Sleep (like stun, but dealing damage interrupts the sleep)

### Positive Effects

Attack effects:

* Attack roll increase
* Attack damage increase (physical/magical)

Skill effects:

* Physical skill roll increase
* Spell skill roll increase

Offensive effects:

* Poison
* Extra damage based on the poison level
* Bleeding
* Extra damage based on the bleeding level

Defensive effects:

* Received attack roll reduction
* Received damage reduction
* Received poison reduction
* Received bleeding reduction
* Shield points (restricted damage reduction)
* Max health increase

Other effects:

* Cure poison
* Convert poison to health
* Convert poison to energy
* Cure bleeding
* Convert bleeding to health
* Convert bleeding to energy
* Max energy increase
* Restore health
* Restore energy

## Bots

### Attack

In some cases, a dice roll is not needed.

A monster will attack an opposing tile if possible.

If there is only a single valid target, that target is selected.

Otherwise, a 1d6 is rolled:

* The dice value represents a target tile
* If that tile is empty, a tile is selected by enumerating the tiles from 1 to 6

If it's impossible to attack the given tile due to the reach or some other reason, a **move action** is used instead to get into the position.

If movement is impossible too, a **guard action** is used instead.

### Movement

1. When **HP is not full** and unit is not in the **back row**, step back to the **back row**
2. Move to the **right**, if possible
3. Move to the **left**, if possible

Otherwise do an **attack action**.
