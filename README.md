# LFU

LFU (least frequently used) - The algorithm is based on a ring data structure

## Description

The simplest implementation of LFU caching (at least frequently used), or rather its modification, on two linked
list/ring.The point is to discard a long-unused item from the cache.The list type structure allows you not to store a
timestamp or an item usage counter.The less frequently an item is used, the further down the list it goes.If the cache
is full, then when adding the "heaviest" element will be deleted.If, when adding an item to the cache, the key is found,
then the value of the item will be updated and it will go to the top of the list.