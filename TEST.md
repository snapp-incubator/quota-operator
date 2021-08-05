## Test plan

1. create a new project:

    ```sh
    oc new-project mytest
    ```

2. confirm default quota and resourcequota are created:

    ```sh
    oc get quota
    oc get resourcequota
    ```

3. try deleting the quota, it should fail with "can not delete default quota"

    ```sh
    oc delete quota default
    ```

4. try updating/replacing the quota, it should fail with "please join your project to a team"

    ```sh
    oc edit quota default
    oc replace -f /tmp/oc-edit-xxx.yaml
    ```

5. join your project to a team without team quota defined, then retry 4, it should fail with "please ask for team quota in cloud support":

    ```
    oc label ns mytest snappcloud.io/team=notExisted
    ```

6. join to a project with team quota defined:

    ```
    oc label ns mytest snappcloud.io/team=snappcloudtest --overwrite
    ```

    then edit quota, it should succeed and also the resourcequota should change accordingly:

    ```
    oc edit quota
    ```
