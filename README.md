# LFU

LFU (least frequently used) - The algorithm is based on a ring data structure

## Description

The simplest implementation of LFU caching (at least frequently used), or rather its modification, on two linked
list/ring.The point is to discard a long-unused item from the cache.The list type structure allows you not to store a
timestamp or an item usage counter.The less frequently an item is used, the further down the list it goes.If the cache
is full, then when adding the "heaviest" element will be deleted.If, when adding an item to the cache, the key is found,
then the value of the item will be updated and it will go to the top of the list.

Example

Initialize the cache of four elements.

```
lfu := NewLFU(4)
```

By itself, it will not be empty - there will be a "zero" element, which in itself is only a boundary for our
data.However, it is important, since we have implemented a two-linked list, the last element will always point to the "
root" element through which we can delete the last element of the list.
<p align='center'><img src="img/1.png?raw=true" alt="root element"></p>


Let's add the first element.In this case, the first element is associated with the root element.

```
lfu.Set("Key_1",1)
```

<p align='center'><img src="img/2.png?raw=true" alt="root element"></p>

Let's add the second element.In this case, a new element is inserted at the beginning of the list, and the first one
goes to the end

```
lfu.Set("Key_2",2)
```

<p align='center'><img src="img/3.png?raw=true" alt="root element"></p>
We will add the third and fourth elements sequentially.As a result, our data structure will look like

```
lfu.Set("Key_3",3)
lfu.Set("Key_4",4)
```

<p align='center'><img src="img/4.png?raw=true" alt="root element"></p>


Since we do not have a counter responsible for marking the number of requests to an element, we will control the "freshness" of the element in a different way.<b>If we request an element using the Get method, we not only return its data, but also move it to the beginning of the queue.</b>The same behavior will occur if we set the value of an already existing key - the element will move to the beginning of the queue, and its value will be updated.For example, let's add an element with an existing key "Key_2".

```
lfu.Set("Key_2","NewValue")
```
Now our structure will look like this

<p align='center'><img src="img/5.png?raw=true" alt="root element"></p>

Now we add the fifth element.Since our size is 4, the cache is already full.When a new item is added, the oldest one will be removed from the cache.
```
lfu.Set("Key_5",5)
```
Our structure takes the following form
<p align='center'><img src="img/6.png?raw=true" alt="root element"></p>

```
lfu.Clear()
```
This method simply clears the cache by deleting all the elements
