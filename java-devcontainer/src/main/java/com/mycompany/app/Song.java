package com.mycompany.app;

import com.datastax.oss.driver.api.mapper.annotations.ClusteringColumn;
import com.datastax.oss.driver.api.mapper.annotations.CqlName;
import com.datastax.oss.driver.api.mapper.annotations.Entity;
import com.datastax.oss.driver.api.mapper.annotations.PartitionKey;

@Entity
@CqlName("songs_by_year")
public class Song {
    @PartitionKey
    @CqlName("year_released")
    private int yearReleased;

    @ClusteringColumn(1)
    @CqlName("artist")
    private String artist;

    @CqlName("song_id")
    private int songId;

    @CqlName("title")
    private String title;
    
    @CqlName("album")
    private String album;

    @CqlName("duration")
    private double duration;

    @CqlName("tempo")
    private double tempo;

    @CqlName("loudness")
    private double loudness;
}
