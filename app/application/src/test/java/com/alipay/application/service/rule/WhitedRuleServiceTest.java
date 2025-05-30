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
package com.alipay.application.service.rule;

import com.alipay.application.service.rule.impl.WhitedRuleServiceImpl;
import com.alipay.application.service.rule.job.ScanService;
import com.alipay.application.share.request.rule.SaveWhitedRuleRequestDTO;
import com.alipay.application.share.request.rule.WhitedRuleConfigDTO;
import com.alipay.common.enums.WhitedRuleOperatorEnum;
import com.alipay.common.enums.WhitedRuleTypeEnum;
import com.alipay.dao.context.UserInfoContext;
import com.alipay.dao.context.UserInfoDTO;
import com.alipay.dao.mapper.RuleMapper;
import com.alipay.dao.mapper.WhitedRuleConfigMapper;
import com.alipay.dao.po.RulePO;
import com.alipay.dao.po.WhitedRuleConfigPO;
import org.junit.Before;
import org.junit.Test;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.Collections;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.when;

/**
 * Date: 2025/4/9
 * Author: lz
 */
public class WhitedRuleServiceTest {

    @InjectMocks
    private WhitedRuleServiceImpl whitedRuleService;

    @Mock
    private WhitedRuleConfigMapper whitedRuleConfigMapper;

    @Mock
    private RuleMapper ruleMapper;

    @Mock
    private ScanService scanService;

    private UserInfoDTO currentUser;

    @Before
    public void setUp() {
        MockitoAnnotations.initMocks(this);
        currentUser = new UserInfoDTO();
        currentUser.setUserId("testUserId");
        currentUser.setTenantId(1111L);
        UserInfoContext.setCurrentUser(currentUser);
    }

    /**
     * [单测用例]测试场景：测试正常情况，插入新的规则
     */
    @Test
    public void testSave_NewRule() {
        SaveWhitedRuleRequestDTO dto = new SaveWhitedRuleRequestDTO();
        dto.setRuleName("testRule");
        dto.setRuleType(WhitedRuleTypeEnum.RULE_ENGINE.name());
        WhitedRuleConfigDTO whitedRuleConfigDTO =  WhitedRuleConfigDTO.builder()
                .id(1)
                .key("key")
                .operator(WhitedRuleOperatorEnum.EQ)
                .value("111")
                .build();
        dto.setRuleConfigList(Collections.singletonList(whitedRuleConfigDTO));
        dto.setEnable(1);
        dto.setRiskRuleCode("testRiskRuleCode");


        when(whitedRuleConfigMapper.insertSelective(any(WhitedRuleConfigPO.class))).thenReturn(1);
        when(ruleMapper.findOne(anyString())).thenReturn(new RulePO());

        int result = whitedRuleService.save(dto);

        assertEquals(1, result);
    }

    /**
     * [单测用例]测试场景：测试正常情况，更新已存在的规则
     */
    @Test
    public void testSave_UpdateRule() {
        SaveWhitedRuleRequestDTO dto = new SaveWhitedRuleRequestDTO();
        dto.setId(1L);
        dto.setRuleName("testRule");
        dto.setRuleType(WhitedRuleTypeEnum.RULE_ENGINE.name());
        WhitedRuleConfigDTO whitedRuleConfigDTO =  WhitedRuleConfigDTO.builder()
                .id(1)
                .key("key")
                .operator(WhitedRuleOperatorEnum.EQ)
                .value("111")
                .build();
        dto.setRuleConfigList(Collections.singletonList(whitedRuleConfigDTO));
        dto.setEnable(1);
        dto.setRiskRuleCode("testRiskRuleCode");

        WhitedRuleConfigPO existingRule = new WhitedRuleConfigPO();
        existingRule.setId(1L);
        existingRule.setLockHolder("testUserId");

        when(whitedRuleConfigMapper.selectByPrimaryKey(anyLong())).thenReturn(existingRule);
        when(whitedRuleConfigMapper.updateByPrimaryKeySelective(any(WhitedRuleConfigPO.class))).thenReturn(1);
        when(ruleMapper.findOne(anyString())).thenReturn(new RulePO());

        int result = whitedRuleService.save(dto);

        assertEquals(1, result);
    }

    /**
     * [单测用例]测试场景：测试异常情况，规则已被其他用户锁定
     */
    @Test(expected = RuntimeException.class)
    public void testSave_RuleLockedByOtherUser() {
        SaveWhitedRuleRequestDTO dto = new SaveWhitedRuleRequestDTO();
        dto.setId(1L);
        dto.setRuleName("testRule");
        dto.setRuleType(WhitedRuleTypeEnum.RULE_ENGINE.name());
        WhitedRuleConfigDTO whitedRuleConfigDTO =  WhitedRuleConfigDTO.builder()
                .id(1)
                .key("key")
                .operator(WhitedRuleOperatorEnum.EQ)
                .value("111")
                .build();
        dto.setRuleConfigList(Collections.singletonList(whitedRuleConfigDTO));
        dto.setEnable(1);
        dto.setRiskRuleCode("testRiskRuleCode");

        WhitedRuleConfigPO existingRule = new WhitedRuleConfigPO();
        existingRule.setId(1L);
        existingRule.setLockHolder("otherUserId");

        when(whitedRuleConfigMapper.selectByPrimaryKey(anyLong())).thenReturn(existingRule);

        whitedRuleService.save(dto);
    }

    /**
     * [单测用例]测试场景：测试异常情况，规则类型不存在
     */
    @Test(expected = RuntimeException.class)
    public void testSave_RuleTypeNotExist() {
        SaveWhitedRuleRequestDTO dto = new SaveWhitedRuleRequestDTO();
        dto.setRuleName("testRule");
        dto.setRuleType("NON_EXISTENT_TYPE");
        dto.setRuleConfigList(Collections.singletonList(new WhitedRuleConfigDTO()));
        dto.setEnable(1);
        dto.setRiskRuleCode("testRiskRuleCode");

        whitedRuleService.save(dto);
    }

    /**
     * [单测用例]测试场景：测试异常情况，用户信息为空
     */
    @Test(expected = RuntimeException.class)
    public void testSave_UserInfoEmpty() {
        SaveWhitedRuleRequestDTO dto = new SaveWhitedRuleRequestDTO();
        dto.setRuleName("testRule");
        dto.setRuleType(WhitedRuleTypeEnum.RULE_ENGINE.name());
        WhitedRuleConfigDTO whitedRuleConfigDTO =  WhitedRuleConfigDTO.builder()
                .id(1)
                .key("key")
                .operator(WhitedRuleOperatorEnum.EQ)
                .value("111")
                .build();
        dto.setRuleConfigList(Collections.singletonList(whitedRuleConfigDTO));
        dto.setEnable(1);
        dto.setRiskRuleCode("testRiskRuleCode");

        when(UserInfoContext.getCurrentUser()).thenReturn(null);

        whitedRuleService.save(dto);
    }



}

