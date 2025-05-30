import LevelTag from '@/components/Common/LevelTag';
import Disposition from '@/components/Disposition';
import {
  obtainGroupTypeTextFromValue,
  obtainPlatformEasyIcon,
} from '@/utils/shared';
import { RightOutlined } from '@ant-design/icons';
import { useIntl, useModel } from '@umijs/max';
import { Button, Divider, Flex, Tag } from 'antd';
import styles from '../index.less';

interface IPolymerizeCard {
  aggregateAsset: API.BaseAggregateAssetInfo;
  cloudAccountId?: string;
}

const PolymerizeCard = (props: IPolymerizeCard) => {
  // Global Info
  const { platformList } = useModel('rule');
  const { groupTypeList } = useModel('asset');
  // Component Props
  const { aggregateAsset, cloudAccountId } = props;
  // Intl API
  const intl = useIntl();
  // Record Info
  const {
    count,
    platform,
    resourceTypeName,
    highLevelRiskCount,
    mediumLevelRiskCount,
    lowLevelRiskCount,
    latestResourceInfo,
    typeFullNameList,
  } = aggregateAsset;

  return (
    <div className={styles['polymerizeCard']}>
      <div className={styles['polymerizeHead']}>
        <Flex align={'center'} style={{ paddingTop: 6 }}>
          <Disposition
            text={resourceTypeName || '-'}
            maxWidth={160}
            rows={1}
            style={{
              color: '#333',
              fontSize: 17,
              fontWeight: 500,
            }}
            placement={'topLeft'}
          />
        </Flex>
        <Disposition
          text={intl.formatMessage(
            {
              id: 'asset.module.text.asset.count',
            },
            {
              count: count,
            },
          )}
          maxWidth={100}
          rows={1}
          style={{
            color: '#FFF',
            fontSize: 14,
          }}
          placement={'topLeft'}
        />
      </div>

      <div className={styles['polymerizeMain']}>
        <div className={styles['riskWrap']}>
          <div className={styles['riskWrapLeft']}>
            <span>{obtainPlatformEasyIcon(platform!, platformList)}</span>
            <Divider type={'vertical'} style={{ margin: '0 8px 0 2px' }} />
            <LevelTag
              level={'HIGH'}
              text={`${intl.formatMessage({
                id: 'common.link.text.high',
              })} ${highLevelRiskCount > 999 ? '999+' : highLevelRiskCount}`}
            />
            <LevelTag
              level={'MEDIUM'}
              text={`${intl.formatMessage({
                id: 'common.link.text.middle',
              })} ${
                mediumLevelRiskCount > 999 ? '999+' : mediumLevelRiskCount
              }`}
            />
            <LevelTag
              level={'LOW'}
              text={`${intl.formatMessage({
                id: 'common.link.text.low',
              })}  ${lowLevelRiskCount > 999 ? '999+' : lowLevelRiskCount}`}
            />
          </div>
          <div className={styles['riskWrapRight']}>
            <Tag className={styles['riskAsset']}>
              {obtainGroupTypeTextFromValue(
                groupTypeList!,
                typeFullNameList?.[0]?.[0],
              )}
            </Tag>
          </div>
        </div>
        <div className={styles['assetMain']}>
          <span className={styles['newTag']}>New</span>
          <div className={styles['assetMainItem']}>
            <span className={styles['assetMainItemLabel']}>
              {intl.formatMessage({
                id: 'cloudAccount.extend.title.asset.name',
              })}
              &nbsp;:
            </span>
            <Disposition
              text={latestResourceInfo?.resourceName || '-'}
              maxWidth={220}
              rows={1}
              style={{
                color: '#333',
                fontSize: 12,
              }}
              placement={'topLeft'}
            />
          </div>
          <div className={styles['assetMainItem']}>
            <span className={styles['assetMainItemLabel']}>
              {intl.formatMessage({
                id: 'cloudAccount.extend.title.asset.id',
              })}
              &nbsp;:
            </span>
            <Disposition
              text={latestResourceInfo?.resourceId || '-'}
              maxWidth={220}
              rows={1}
              style={{
                color: '#333',
                fontSize: 12,
              }}
              placement={'topLeft'}
            />
          </div>
          <div className={styles['assetMainItem']}>
            <span className={styles['assetMainItemLabel']}>
              {intl.formatMessage({
                id: 'cloudAccount.extend.title.updateTime',
              })}
              &nbsp;:
            </span>
            <Disposition
              text={latestResourceInfo?.gmtModified || '-'}
              maxWidth={220}
              rows={1}
              style={{
                color: '#333',
                fontSize: 12,
              }}
              placement={'topLeft'}
            />
          </div>
          <div className={styles['assetMainItem']}>
            <span className={styles['assetMainItemLabel']}>
              {intl.formatMessage({
                id: 'asset.module.input.text.ip',
              })}
              &nbsp;:
            </span>
            <Disposition
              text={latestResourceInfo?.address || '-'}
              maxWidth={220}
              rows={1}
              style={{
                color: '#333',
                fontSize: 12,
              }}
              placement={'topLeft'}
            />
          </div>
        </div>
        <div className={styles['assetView']}>
          <Button
            href={
              cloudAccountId
                ? `/assetManagement/assetList?platform=${platform}&resourceGroupType=${typeFullNameList?.[0]?.[0]}&resourceType=${typeFullNameList?.[0]?.[1]}&cloudAccountId=${cloudAccountId}`
                : `/assetManagement/assetList?platform=${platform}&resourceGroupType=${typeFullNameList?.[0]?.[0]}&resourceType=${typeFullNameList?.[0]?.[1]}`
            }
            type={'link'}
            style={{ fontSize: 14, gap: 4 }}
          >
            {intl.formatMessage({
              id: 'common.button.text.viewDetail',
            })}
            <RightOutlined />
          </Button>
        </div>
      </div>
    </div>
  );
};
export default PolymerizeCard;
