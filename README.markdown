# Chat de go seguro

Tenemos un chat donde pueden loguear varios usurios y tener una conversación entre ellos. Se quiere aplicar seguridad al chat.

## Pequeñas modificaciones del chat

En la vista mostramos ahora tambien el nombre del usuario que envia el mensaje; ya que si entraran varios urusarios la imagen solo no seria identificatiba de quien envia el mensaje.

## Seguridad

Hemos añadido el cifrador RC4 para cifrar los mensajes, como podemos ver el metodo push ahora enviamos el mensaje cifrado.

Para luego desde el pull descifrar ese mensaje.

## Autenticacion

Ahora se puede leer desde fichero de texto y autenticar contra un fichero de texto donde tendremos los distintos usuarios y sus contraseñas.

Antiguamente el hash que utilizabamos para las contrasenyas era md5 pero ahora lo hemos remplazado por sha512 debido a la poca seguridad de md5.