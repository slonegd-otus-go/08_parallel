# Параллельное исполнение
Сделать функцию для параллельного выполнения N заданий.
Принимает на вход слайс с заданиями 
```golang
[]func()error
```
Число заданий которые можно выполнять параллельно `N` и максимальное число ошибок после которого нужно приостановить обработку. 

Учесть что задания могу выполняться разное время. 

# reviev
Предлагаю wait group использовать только для воркеров. Перед стартом каждого воркера https://github.com/slonegd-otus-go/08_parallel/blob/7f55cf7bd92bec12924fb2510d385e70bbaf6682/parallel.go#L14 делать waitgroup.Add(1), а в самой функции worker делать defer waitgroup.Done() Таким образом, можно не учитывать waitgroup для Execute и функций task. 

Бонусом и от этой конструкции удастся избавиться https://github.com/slonegd-otus-go/08_parallel/blob/7f55cf7bd92bec12924fb2510d385e70bbaf6682/parallel.go#L33-L39 переносом waitgroup.Wait() в конец функции Execute 

Можно пройтись range по задачам https://github.com/slonegd-otus-go/08_parallel/blob/7f55cf7bd92bec12924fb2510d385e70bbaf6682/parallel.go#L42

