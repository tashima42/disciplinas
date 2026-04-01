require(ggmap)

load('../all_bus_info')
all_bus_info

readRenviron('.env')
STADIA_API_KEY <- Sys.getenv("STADIA_API_KEY")

register_stadiamaps(STADIA_API_KEY, write = TRUE)

df_lines <- data.frame(index = numeric(length(all_bus_info)),
                       code = character(length(all_bus_info)),
                       name = character(length(all_bus_info)),
                       cat = character(length(all_bus_info)),
                       stringsAsFactors = FALSE)

for (i in 1:length(all_bus_info)) {
  df_lines$index[i] <- i
  df_lines$code[i] <- all_bus_info[[i]][[1]]
  df_lines$name[i] <- all_bus_info[[i]][[2]]
  df_lines$cat[i] <- all_bus_info[[i]][[3]]
}

df_all_stops <- data.frame(num = character(0),
                           lat = numeric(0),
                           lon = numeric(0),
                           group = character(0),
                           cat = character(0),
                           stringsAsFactors = FALSE)

for (i in 1:length(all_bus_info)) {
  num_of_directions <- length(all_bus_info[[i]][[4]])
  if (num_of_directions == 0)
    next
  
  a_line_code <- all_bus_info[[i]][[1]]
  
  for (j in 1:num_of_directions) {
    total_stops <- nrow(all_bus_info[[i]][[4]][[j]][[2]])
    df_sub <- data.frame(num=character(total_stops),
                         lat = numeric(total_stops),
                         lon = numeric(total_stops),
                         group = character(total_stops),
                         cat = character(total_stops),
                         stringsAsFactors = FALSE)
    df_sub$num <- all_bus_info[[i]][[4]][[j]][[2]]$NUM
    df_sub$lat <- all_bus_info[[i]][[4]][[j]][[2]]$LAT
    df_sub$lon <- all_bus_info[[i]][[4]][[j]][[2]]$LON
    df_sub$group <- all_bus_info[[i]][[4]][[j]][[2]]$GRUPO
    df_sub$cat <- df_lines[df_lines$code==a_line_code,]$cat
    
    df_all_stops <- rbind(df_all_stops, df_sub)
  }
}

df_all_stops <- df_all_stops[!duplicated(df_all_stops$num), ]

city_center = c(mean(df_all_stops$lon), mean(df_all_stops$lat))

curitiba <- c(left=-49.38039, top=-25.34738, right=-49.17349, bottom=-25.63794)
my_curitiba_v11 <- get_stadiamap(curitiba, zoom = 11)
ggmap(my_curitiba_v11, extent = 'device') + geom_point(aes(x = lon, y = lat, color = factor(cat)),
                                                       data = df_all_stops)
