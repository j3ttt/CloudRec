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
package com.alipay.dao.mapper;

import com.alipay.dao.dto.TenantDTO;
import com.alipay.dao.dto.UserDTO;
import com.alipay.dao.po.UserPO;
import org.apache.ibatis.annotations.Param;

import java.util.List;

public interface UserMapper {
    int deleteByPrimaryKey(Long id);

    int insertSelective(UserPO record);

    UserPO selectByPrimaryKey(Long id);

    int updateByPrimaryKeySelective(UserPO record);

    UserPO findOne(String userId);

    int findMemberCount(TenantDTO tenantDTO);

    List<UserPO> findMemberList(TenantDTO tenantDTO);

    int findCount(UserDTO userDTO);

    List<UserPO> findList(UserDTO userDTO);

    UserPO find(@Param("userId") String userId, @Param("password") String password);

    UserPO findByUserName(@Param("userName") String userName);

}