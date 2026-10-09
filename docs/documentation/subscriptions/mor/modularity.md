
To design the modular architecture for the mor we must take into account
- A mor supplier is used such as paddle or polar but not both
- Each supplier handles differently: paddle you can create separate prices polar no 
- then just an adapter would not be enough

Who will be in charge of bringing the plans why?
- /plans only brings the plans but the price is not there at all and we cannot call 
- some provider because it would generate plan dependencies and it should be the other way around
- then to bring the plans it will be within each provider to get plans from there we join it with the data
  from each source
- Of course, before enabling the supplier, all the data must be synchronized from the suppliers.
  and also with the collection
- How do you use one or the other then?
- For example, inside paddle we can use Google Pay from your dashboar you can configure Google
- Since they are subcriptions I will be consulting the subpiptions at paddle, generally there are not many records
  so I'll just be checking


- doble habilitar debes tener relacionado los productos y validar que la relacion plan con plan de paddle
 no sea  nulo por que  en tiempo de ejecucion alguien puede crear el producto, y estara nil