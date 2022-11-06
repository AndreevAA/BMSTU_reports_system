FROM 'byjg/nginx-extras'

wget 'http://nginx.org/download/nginx-1.17.8.tar.gz'
tar -xzvf nginx-1.17.8.tar.gz
cd nginx-1.17.8/

# Here we assume you would install you nginx under /opt/nginx/.
./configure --prefix=/opt/nginx \
 --add-module=/path/to/headers-more-nginx-module

make
make install

# Set working directory
WORKDIR /backend

# Copy over static assets
COPY static static/

# Copy over entrypoint
COPY nginx-entrypoint.sh /usr/local/bin/nginx-entrypoint.sh

# Copy Nginx config template
RUN rm /etc/nginx/conf.d/default.conf
COPY nginx.conf /etc/nginx/conf.d/nginx.conf

# Nginx init
RUN ["chmod", "+x", "/usr/local/bin/nginx-entrypoint.sh"]
ENTRYPOINT ["/usr/local/bin/nginx-entrypoint.sh"]