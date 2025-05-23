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
package com.alipay.dao.dto;

import com.fasterxml.jackson.annotation.JsonFormat;
import lombok.Getter;
import lombok.Setter;

import java.util.Date;
import java.util.List;

@Getter
@Setter
public class ResourceAggByInstanceTypeDTO {

    private String platform;

    private Integer count;

    private String resourceType;

    private String resourceTypeName;

    private List<List<String>> typeFullNameList;

    private Integer highLevelRiskCount;

    private Integer mediumLevelRiskCount;

    private Integer lowLevelRiskCount;

    private LatestResourceInfo latestResourceInfo;

    @Getter
    @Setter
    public static class LatestResourceInfo {
        private String resourceId;

        private String resourceName;

        @JsonFormat(pattern = "yyyy-MM-dd HH:mm:ss", timezone = "GMT+8")
        private Date gmtModified;

        private String address;
    }
}
