# FizzBuzz

### Endpoints
#### /fizzbuzz
* Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.
* Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

* Usage: POST a JSON object with the five parameters as named above. 

```
{   
    "limit": 100,
    "int1": 3,
    "int2": 5,
    "string1": "Fizz",
    "string2": "Buzz"
 }
```
#### /statistics
* Returns the parameters corresponding to the most used request, as well as the number of hits for this request
* Usage: GET could return the following response  
```
[
    {
        "parameters": {
            "limit": 100,
            "int1": 3,
            "int2": 5,
            "string1": "Fizz",
            "string2": "Buzz"
        },
        "count": 4
    }
]

```