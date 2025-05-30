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
package com.alipay.application.service.risk.domain.repo;


import com.alibaba.fastjson.JSON;
import com.alipay.dao.dto.Subscription;
import com.alipay.dao.converter.Converter;
import com.alipay.dao.po.SubscriptionPO;
import org.springframework.beans.BeanUtils;
import org.springframework.stereotype.Component;

/*
 *@title SubscriptionConverter
 *@description
 *@author jietian
 *@version 1.0
 *@create 2025/1/17 22:37
 */
@Component
public class SubscriptionConverter implements Converter<Subscription, SubscriptionPO> {

    @Override
    public SubscriptionPO toPo(Subscription subscription) {
        SubscriptionPO subscriptionPO = new SubscriptionPO();
        BeanUtils.copyProperties(subscription, subscriptionPO);
        subscriptionPO.setActionList(JSON.toJSONString(subscription.getActionList()));
        return subscriptionPO;
    }

    @Override
    public Subscription toEntity(SubscriptionPO subscriptionPO) {
        Subscription subscription = new Subscription();
        BeanUtils.copyProperties(subscriptionPO, subscription);
        subscription.setActionList(JSON.parseArray(subscriptionPO.getActionList(), Subscription.Action.class));
        return subscription;
    }
}
