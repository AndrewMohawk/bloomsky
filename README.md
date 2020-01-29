# BloomSky
Few simple scripts to accompany https://www.andrewmohawk.com/2020/01/28/picking-apart-an-iot-camera-bloomsky/ that target http://www.bloomsky.com/ devices

# Accounts
If you need an account simply download the app for your mobile device and register. The registration is not email validated so feel free to use any email address

# Building
Just compile with go on whatever you are running:

 `go build bloomsearch.go`
 
 `go build bloompost.go`
 
# BloomSearch
Lets you search the same as within the app interface to return the up to 20 results based on a keyword, while it is possible and described within the blog post that you can enumerate all devices this script will NOT do that for you.

```
bloomsearch <username> <password> <searchterm>
                        username - Your bloomsky username
                        password - Your bloomsky password
                        searchTerm - Bloomsky you are searching for
```

Search results will vary wildly and honestly, I am not sure why, results are consistent with what is in the mobile application, here is an example searching for my device 'AndrewMohawk':

```
./bloomsearch.exe bloomtest1@andrewmohawk.com mycoolpassword AndrewMohawk
Got Auth Token, searching for AndrewMohawk
Making API request to https://bskybackend.bloomsky.com/api/skydevice/find_me/?name=AndrewMohawk
Token mycoolToken
Device Name                              | Device ID                 | Street Name 
-----------                              | ---------                 | ----------- 
AndrewMohawk                             | 94A1A2733B1A              | Bush Street
St Andrew UMC                            | ************              | North 50 West
Dan  Andrews West                        | ************              | Walnut Branch Lane
Mohawk Trail Regional School             | ************              | Ashfield Rd
Mark @ Red Hawk                          | ************              | Thistle Belle Court
andrelix                                 | ************              | Faith Ave
Downtown211.com  danandrews.com          | ************              | Market Sq
Søndre Led                               | ************              | Søndre Led
Hawks Nest Bloom                         | ************              | E Hawks Nest Pl
Andrei's bloomsky                        | ************              | Foothill Oaks Terrace
Lily Landry                              | ************              | Lily Landry Ct
Whitehawk                                | ************              | Woodson Dr
Fishhawk Ranch                           | ************              | Bridgewalk Dr
Tomahawk Trail                           | ************              | Tomahawk Trail
SOHA STL                                 | ************              | Lansdowne Ave
Hawk Bait Weather Station                | ************              | Horizon Circle
Alexandria                               | ************              | Manchester Park Cir
J and J W                                | ************              | Granite Lane
Weatherhawk                              | ************              | E 1075 N
Hawkshead weather                        | ************              | Private Street
```

# BloomPost
No Authentication/Authorisation is required, simply specify the DeviceID as well as the path to the image and you can upload an image, static values are set for all the fields, but feel free to adjust as needed 
```
bloompost <DeviceID> <imagePath>
         DeviceID - Device ID of your Bloomsky camera, use bloomsearch to find this or look at the traffic from your device
         imagePath - full path of the image to upload
```
Example:
```
 .\bloompost.exe 94A1A2733B1A ..\wombat.jpg
Updated Bloomsky!
```


# Bloomproxy
Hacky code to demonstrate making a somewhat proxy for this scenario.
