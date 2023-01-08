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
