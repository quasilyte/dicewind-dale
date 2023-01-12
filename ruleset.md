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

* Lockpicking
* Search
* Trap disarming
* Appraise
* Alchemy
* Item identification

**Lockpicking** allows a character to open a locked object with 100% chance. Without such talent, there is only ~16.6% of success (6 should be rolled).

**Search** allows to detect a hidden object or room, opening more possibilities for exploration.

**Trap disarming** allows a character to safely disarm any trap. Some random events may have trap-related moments. Some hidden rooms may be guarded by a trapped door. In addition to that, a treasure chest can be protected by a trap in addition to being locked.

**Appraise** improves the prices in the shops. You can sell items for more gold and buy new equipment for lower prices.

**Alchemy** allows potion mixing to achieve unique results. It also improves the potion effects when they're consumed by the character with this talent.

**Item identification** allows a party to identify any new item at once, anywhere. Without this talent, items can be identified only in dedicated places like shops.

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
