Puntos importantes al diseñar la db

1. Verificar que tablas tendrán created_at y updated_at
2. Como se manejara las eliminaciones eliminando en disco (aqui entraria el tema de eliminacion en cascada claro otra forma es usando transacciones para asegurarnos que sean solo las tablas que queremos), poniendo en otra tabla, o creo mas simple es 
   añadiendo un campo deleted_at a la tabla y en el servicio validar, si fue eliminado,


Migraciones
1. ¿Cuándo usar migraciones?
2. ¿Cómo ejecutarlas segun ambientes de prod y desarollo? y que no afecten a prod 