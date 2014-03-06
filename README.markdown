# Chat de go seguro

Tenemos un chat donde pueden loguear varios usurios y tener una conversación entre ellos. Se quiere aplicar seguridad al chat.

## Pequeñas modificaciones del chat

En la vista mostramos ahora tambien el nombre del usuario que envia el mensaje; ya que si entraran varios urusarios la imagen solo no seria identificatiba de quien envia el mensaje.

## Seguridad

Hemos añadido el cifrador RC4 para cifrar los mensajes, como podemos ver el metodo push ahora enviamos el mensaje cifrado.

Para luego desde el pull descifrar ese mensaje.