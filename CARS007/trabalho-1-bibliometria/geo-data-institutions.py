from dotenv import dotenv_values
# from geopy.geocoders import Nominatim
from geopy.geocoders import ArcGIS
import json
import time

config = dotenv_values(".env")

# geolocator = Nominatim(user_agent="UniversityLocator 0.0.1")

geolocator = ArcGIS(user_agent="UniversityLocator 0.0.1")

institutions_file = config.get("INSTITUTIONS")
geo_file = config.get("INSTITUTIONS_GEO_DATA")

# institutions = []
# with open(institutions_file, "r") as inst_file:
#     institutions = json.load(inst_file)

geo_data = {}

with open(geo_file, "r") as data:
  geo_data = json.load(data)

institutions = geo_data["missing"]

for inst in institutions:
  if inst in geo_data["found"]:
    print(f"Already geocoded {inst}, skipping")
    continue
  with open(geo_file, "w") as outfile:
    try:
      location = geolocator.geocode(inst)
      if location:
        geo_data["found"][inst] = {
          "latitude": location.latitude,
          "longitude": location.longitude,
          "address": location.address,
        }
        print(f"Geocoded {inst}: {location.latitude}, {location.longitude}")
      else:
        geo_data["not_found"].append(inst)
        print(f"Could not geocode {inst}")
    except Exception as e:
      geo_data["errors"].append(inst)
      print(f"Error geocoding {inst}: {e}")
    json.dump(geo_data, outfile, indent=2)
  time.sleep(1)