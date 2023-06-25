## Task
Prepare an architecture document specifying an Instagram-like application allowing users to
store media files and that stores proof of the existence of media files on the Ethereum chain.

### Acceptance Criteria
* Users can upload/download/view media files such as pictures and videos.
* Users can check the existence of the validity of the media files by checking the proof of
existence on-chain.
* Users can perform searches based on media file titles.


## Solution

### System requirements
* Typically, such systems handle more read requests than write requests. Therefore, our system should be optimized for read requests.
* System should be scalable and fault-tolerant, ensuring high availability and the ability to handle increasing load.
* Data integrity is crucial. Once media content is uploaded and the user receives confirmation, it must be securely stored and not lost.
* The system should provide low latency for read requests, ensuring fast retrieval of media content.

### Calculations
* Let's assume that we have 100_000_000 users. 
* Let's assume that about 500_000 media files will be uploaded daily.
* Let's assume that average size of media file is 1Mb.

Based on this assumptions we can calculate that we will have about 15_000_000 media files in 30 days. 

If we store all the media files in one place, we will accumulate around 15TB of data within that period, resulting in approximately 180TB per year.

While this amount of data may not pose a significant challenge for modern hardware, it's important to consider the high volume of read requests we expect. Therefore, we should focus on optimizing our system for efficient read operations.

By designing the system with a focus on read performance, we can ensure that users can quickly access and retrieve media files, providing a seamless experience with low latency.

#### Media meta info
* Media meta info associated with media file only.

### High-level design

#### File upload
* Uploading files usually takes more time compared to downloading.
* The file upload process can consume all available server connections.

Based on these points, it is recommended to use separate servers for file uploading and downloading. This approach has several benefits:
* Prevents exhausting the available connections on the server during file uploads.
* Allows for independent scaling of the uploading and downloading components, addressing potential bottlenecks.
* Provides flexibility to optimize each server for its specific task, such as fine-tuning network settings for efficient uploading or caching strategies for fast downloading.

By segregating the upload and download processes onto separate servers, we can ensure better performance, scalability, and fault tolerance in our system.

#### File storage
* To avoid data loss, we employ a multiple-copy storage approach. If one storage location becomes unavailable, we replace it with another location where copies are stored.
* The same principle applies to other system components. It is always important to maintain multiple replicas to handle failures. Replication helps eliminate single points of failure in the system.
* If a service requires a copy, it can be kept in the background but not actively serving traffic until the primary copy fails.

Backup systems or redundant functionality enable the system to recover in crisis situations.
In the event of a failure, the system can switch to a functioning replica.

### System components
* **AppSC**: smart contract that stores information about all the users in the system. 
  * each user has a dedicated smartcontract which stores information about all the files uploaded by the user.
* **UserSC**: user's smart contract that stores information about all the files uploaded by the user.
  * keep track of proof for each file from *storageCluster*
  * keep list of uploaded file names
* **Storage clusters**: each cluster stores metadata about images and maintains a connection to its dedicated storage location where the image is stored.
  * Metadata database (meta db): The meta db stores image metadata properties. Initially, it can be located on the same database as the storage clusters and scaled as needed based on demand.
    * metainfo we need at first time *(76 bytes per image, 1140Mb for 1 month)*:
      * image id *(32 bytes)*
      * latitute *(4 bytes)*
      * longitude *(4 butes)*
      * title *(32 bytes)*
      * creation date *(4 bytes)*
  * Storage: The storage component is responsible for handling the reading and writing of image files on demand.
* **TopologySC**: smart contract that stores information about all the storage clusters in the system. Each cluster in the system has its own unique public-private key pair, which is used for the verification of file storage operations.
  * acts as a centralized registry for all the storage clusters, maintaining a record of their public keys and other relevant information.
* **The Consistent Hasher module**: responsible for ensuring the equal distribution of image traffic across the storage clusters in the system. 
  * use consistent hash algorithm to assign each storage cluster to a position on a virtual circle. 
  * allows for quick and efficient scaling or replacement of storage clusters within the system.
  * key calculation based on the combination of the "imageHash" and a specific "key" value - predictable and consistent distribution of image traffic among the storage clusters. This ensures that the workload is evenly distributed across the clusters, preventing hotspots and optimizing the overall system performance.
* **KeyGenerator**: service with its own dedicated database and dedicated nodes. 
  * During the system's initialization, a certain number of sequential numbers are generated, for example (let's assume 1_000_000 an initial generation). Additionally, an extra 1_000 keys are generated in the background.
  * Each node in service loads a portion of these numbers from the database into its internal memory and marks them as used. When a request for a key is received, the service provides the first available number from its internal memory and removes it from the memory.
  * In the event of a service failure, there are other numbers available in the sequential sequence, so there is no need for a backup.
* **Observer**: service which observe storage clusters state and change **TopologySC** if it need. On new cluster join into system - consistent hash circle should be rebalanced and **TopologySC** should be updated.
* **Onboarder**: service is responsible for registering new users and deploying personalized smart contracts for each user. It also stores the associated contract address for each user in the database.

### Detailed design

#### User registration flow
1. User registers in the system via wallet.
2. System ask user to deposit some gas to *AppSC* for smart contract deployment and future-time proof upload.
3. User deposits gas to *AppSC* and system deploys new smart contract for user. Address of new smart contract is stored in *AppSC*.

#### File upload flow
1. User make photo and start upload it to the system.
2. Service make photo compression and generate smallers versions of photo (thumbnails).
3. After compression and generation of thumbnails, system show's result to user and ask him to confirm upload.
4. User confirm and sign upload request. Each generated thumbnail has own hash and own user's sign.
5. Gateway service start upload process:
   1. Get user's smart contract address from **AppSC**.
   2. Get new key from **KeyGenerator**.
   3. Calculate hash of image and get storage cluster address from **Consistent Hasher**.
   4. Upload image to storage cluster.
   5. Wait for confirmation from storage cluster.
   6. **Storage cluster** save files and sign it in confirmation.
   6. If confirmation received - update **UserSC** and return result to user.

HumanReadable structure stored in used in **UserSC**. It can be optimized for use less memory and gas with encoding structs and data:
```json
{
  "fileName": "image.jpg",
  "objects": [
    {
      "type": "original",
      "hash": "0x1234567890",
      "proofs": [
        {
          "storageCluster": "0x1234567890",
          "signature": "0x1234567890"
        },
        {
            "storageCluster": "0x1234567890",
            "signature": "0x1234567890"
        }
      ]
    },
    {
        "type": "thumbnail",
        "hash": "0x1234567890",
        "proofs": [
            {
                "storageCluster": "0x1234567890",
                "signature": "0x1234567890"
            },
          {
              "storageCluster": "0x1234567890",
              "signature": "0x1234567890"
          }
        ]
    }
  ],
}
```

#### File download flow
1. User request file from system.
2. Gateway service get user's smart contract address from **AppSC**.
3. Gateway service get file from **UserSC**.
4. Gateway service get storage cluster address from **Consistent Hasher**.
5. Gateway service request file from storage cluster.
6. Gateway service return file to user.

#### File deletion flow
1. User request file deletion from system.
2. Gateway service get user's smart contract address from **AppSC**.
3. Gateway service get file from **UserSC**.
4. Gateway service get storage cluster address from **Consistent Hasher**.
5. Gateway service request file deletion from storage cluster.
6. Gateway service return result to user.

#### File search flow
1. User request file search from system.
2. Gateway service get user's smart contract address from **AppSC**.
3. Gateway service get list of files **UserSC**.
4. If file found - get file info from **UserSC**.
4. Based of file info - get storage cluster address from **Consistent Hasher**.
5. Gateway service request file from storage cluster.
6. Gateway service return file to user.

#### File storage proof check flow
1. User request file storage proof check from system.
2. Gateway service get user's smart contract address from **AppSC**.
3. Gateway service get file from **UserSC**.
4. Gateway load storage proofs from file metadata. 
5. Gateway service get list of storage clusters from **TopologySC** and verify signatures.