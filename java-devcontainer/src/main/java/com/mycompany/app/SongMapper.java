package com.mycompany.app;

import com.datastax.oss.driver.api.mapper.annotations.DaoFactory;
import com.datastax.oss.driver.api.mapper.annotations.Mapper;

@Mapper
public interface SongMapper {
    @DaoFactory
    SongDao songDao();
}
