package com.mycompany.app;

import java.net.InetSocketAddress;
import java.time.Duration;

import com.datastax.oss.driver.api.core.CqlIdentifier;
import com.datastax.oss.driver.api.core.CqlSession;
import com.datastax.oss.driver.api.core.DefaultConsistencyLevel;
import com.datastax.oss.driver.api.core.config.DefaultDriverOption;
import com.datastax.oss.driver.api.core.config.DriverConfigLoader;
import com.datastax.oss.driver.api.core.cql.*;

public class App {

    public static void main(String[] args) {        
        // Config
        DriverConfigLoader loader =
                DriverConfigLoader.programmaticBuilder()
                        .withString(DefaultDriverOption.REQUEST_CONSISTENCY, DefaultConsistencyLevel.LOCAL_QUORUM.name())
                        .withDuration(DefaultDriverOption.REQUEST_TIMEOUT, Duration.ofSeconds(2))
                        .build();
        
        // Build the session
        CqlSession session = CqlSession.builder()
            .addContactPoint(new InetSocketAddress("scylla-node1", 9042))
            .addContactPoint(new InetSocketAddress("scylla-node2", 9042))
            .addContactPoint(new InetSocketAddress("scylla-node3", 9042))
            .withLocalDatacenter("DC1")
            .withKeyspace(CqlIdentifier.fromCql("songs"))
            .withConfigLoader(loader)
            .build();
        
        // Quick example
        ResultSet rs = session.execute("SELECT * FROM songs_by_year");

        for (Row row : rs) {
            System.out.println(row.getString("title"));
        }        
    }
}
