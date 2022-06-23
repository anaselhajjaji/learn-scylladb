package com.mycompany.app;

import com.datastax.oss.driver.api.mapper.annotations.Dao;
import com.datastax.oss.driver.api.mapper.annotations.Select;

@Dao
public interface SongDao {
    @Select
    Song findSongs();
}
