/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package com.alipay.bootstrap;

import lombok.extern.slf4j.Slf4j;
import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@MapperScan("com.alipay.dao.mapper")
@ComponentScan(basePackages = {"com.alipay.application", "com.alipay.api",
        "com.alipay.dao"})
@Slf4j
@SpringBootApplication
public class BootstrapApplication {

    public static void main(String[] args) {
        System.setProperty("jdk.tls.maxHandshakeMessageSize", "51200");
        SpringApplication.run(BootstrapApplication.class, args);
        log.info("Application start success...");
    }

}
