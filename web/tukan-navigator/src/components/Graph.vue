<template>
  <v-container fluid class="fill-height">
    <v-row class="fill-height">
      <v-col class="fill-height">
        <!-- SVG container with relative positioning -->
        <div class="svg-container">
          <!-- SVG graph with absolute positioning -->
          <svg
            ref="svg"
            class="graph-svg"
            @mousedown="startInteraction"
            @touchstart="startInteraction"
            @mouseup="stopInteraction"
            @mouseleave="stopInteraction"
            @touchend="stopInteraction"
            @mousemove="handleMouseOver"
            @touchmove="handleTouchMove"
            @click="handleSvgClick"
          >
            <g :transform="`translate(${panX}, ${panY})`">
              <!-- Define grid pattern -->
              <defs>
                <pattern
                  id="gridPattern"
                  width="40"
                  height="40"
                  patternUnits="userSpaceOnUse"
                >
                  <rect width="40" height="40" fill="#f0f0f0" />
                  <path
                    d="M 40 0 L 0 0 0 40"
                    fill="none"
                    stroke="#cccccc"
                    stroke-width="1"
                  />
                </pattern>
              </defs>
            </g>
              
            <g :transform="`translate(${panX}, ${panY})`">
              <!-- Apply grid pattern as background -->
              <rect
                height="100000"
                width="100000"
                x="-50000"
                y="-50000"
                fill="url(#gridPattern)"
              />
            </g>
            
            <!-- Existing SVG content -->
            <g :transform="`translate(${panX}, ${panY})`">
              <!-- Render existing edges -->
                <Edge
                v-for="edge in validEdges"
                :key="'edge-' + edge.id"
                :from="findNode(edge.from)"
                :to="findNode(edge.to)"
                :edge="edge"
                :highlighted="(selectedEdgeId === edge.id || highlightedEdgeId === edge.id) && (mode === 'edit' || mode === 'remove')"
                @select="handleEdgeSelect(edge)"
                />
              />
              <!-- Render existing nodes -->
                <Node
                v-for="node in nodes"
                :key="'node-' + node.id"
                :id="node.id"
                :x="node.x"
                :y="node.y"
                :selected="node.id === selectedNodeId || node.id === edgeStartNode"
                :highlighted="node.id === hoveredNodeId || node.id === selectedNodeIdInEditMode"
                :mode="mode"
                @select="selectNode"
                @mousedown.stop="startNodeDrag($event, node.id)"
                @touchstart.stop="startNodeDrag($event, node.id)"
                />
              />
            </g>

            <!-- Ghost element for new edge -->
            <line
              v-if="mode === 'addEdge' && edgeStartNode && hoveredNodeId"
              :x1="findNode(edgeStartNode).x + panX"
              :y1="findNode(edgeStartNode).y + panY"
              :x2="findNode(hoveredNodeId).x + panX"
              :y2="findNode(hoveredNodeId).y + panY"
              stroke="black"
              opacity="0.5"
              z-index="1"
              stroke-dasharray="5,5"
              pointer-events="none"
            />

            <!-- Ghost element for new node -->
            <circle
              v-if="mode === 'add' && addNodeHovered"
              :cx="addNodeHovered.x + panX"
              :cy="addNodeHovered.y + panY"
              r="20"
              fill="transparent"
              stroke="black"
              opacity="0.5"
              z-index="1"
              stroke-dasharray="5,5"
              pointer-events="none"
            />

            <Toucan :x="toucanX" :y="toucanY" />
          </svg>

          <!-- Button group outside SVG for z-index stacking -->
          <div class="button-group top">
            <!-- Buttons with icons and gaps -->
            <v-btn
              :class="{ 'v-btn--active': mode === 'pan' }"
              @mousedown.stop="setMode('pan')"
              @touchstart.stop="setMode('pan')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-cursor-move</v-icon> Move
            </v-btn>

            <v-btn
              :class="{ 'v-btn--active': mode === 'drag' }"
              @mousedown.stop="setMode('drag')"
              @touchstart.stop="setMode('drag')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-drag-variant</v-icon> Drag
            </v-btn>

            <v-btn
              :class="{ 'v-btn--active': mode === 'add' }"
              @mousedown.stop="setMode('add')"
              @touchstart.stop="setMode('add')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-vector-circle</v-icon> Add Node
            </v-btn>

            <v-btn
              :class="{ 'v-btn--active': mode === 'addEdge' }"
              @mousedown.stop="setMode('addEdge')"
              @touchstart.stop="setMode('addEdge')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-vector-line</v-icon> Add Edge
            </v-btn>
          </div>
          
          <div class="button-group bottom">
            <v-btn
              :class="{ 'v-btn--active': mode === 'edit' }"
              @mousedown.stop="setMode('edit')"
              @touchstart.stop="setMode('edit')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-pencil</v-icon> Edit
            </v-btn>

            <v-btn
              :class="{ 'v-btn--active': mode === 'remove' }"
              @mousedown.stop="setMode('remove')"
              @touchstart.stop="setMode('remove')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-delete</v-icon> Remove
            </v-btn>

            <v-btn
              @mousedown.stop="panToNode('S')"
              @touchstart.stop="panToNode('S')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-page-first</v-icon> Go to start
            </v-btn>

            <v-btn
              @mousedown.stop="panToNode('P')"
              @touchstart.stop="panToNode('P')"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
            >
              <v-icon>mdi-page-last</v-icon> Go to end
            </v-btn>

            <v-btn
              :class="animationError ? 'error' : (animationRunning ? 'warning' : 'success')"
              @mousedown.stop="toggleAnimation"
              @touchstart.stop="toggleAnimation"
              :size="isMobile ? 'x-small' : 'large'"
              density="comfortable"
              :color="animationRunning ? 'error' : (animationError ? 'warning' : 'success')"
            >
              <v-icon>{{ animationRunning ? 'mdi-stop' : (animationError ? 'mdi-replay' : 'mdi-play') }}</v-icon>
              {{ animationRunning ? 'Stop' : (animationError ? 'Restart' : 'Start') }}
            </v-btn>
          </div>
        </div>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { throttle } from 'lodash';
import Node from './Node.vue';
import Edge from './Edge.vue';
import Toucan from './Toucan.vue';

export default {
  components: {
    Node,
    Edge,
    Toucan
  },
  data() {
    return {
      nodes: [
        { id: 'S', x: 50, y: 50 },
        { id: 'P', x: 250, y: 50 },
      ],
      edges: [],
      selectedNodeId: null,
      hoveredNodeId: null,
      selectedEdgeId: null,
      highlightedEdgeId: null,
      selectedNodeIdInEditMode: null,
      draggingNodeId: null,
      isPanning: false,
      offsetX: 0,
      offsetY: 0,
      panX: 500,
      panY: 500,
      panStartX: 0,
      panStartY: 0,
      toucanX: 0,
      toucanY: 0,
      width: 800,
      height: 600,
      mode: 'pan',
      creatingEdge: null,
      edgeStartNode: null,
      addNodeHovered: null,
      animationRunning: false,
      animationError: false,
      animatedPath: [],
      pathIndex: 0,
      isMobile: true,
      isPathSet: false,
    };
  },
  computed: {
    validEdges() {
      return this.edges.filter(edge => this.findNode(edge.from) && this.findNode(edge.to));
    },
    centerX() {
      return -this.panX;
    },
    centerY() {
      return -this.panY;
    }
  },
  mounted() {
    this.throttledMouseMove = throttle(this.onMouseMove, 16);
    this.detectMobile();
  
    window.addEventListener('resize', this.detectMobile);
    window.addEventListener('mousemove', this.throttledMouseMove);
    window.addEventListener('mouseup', this.stopInteraction);

    // Touch event listeners
    window.addEventListener('touchmove', this.onTouchMove, { passive: false });
    window.addEventListener('touchend', this.stopInteraction);
    window.addEventListener('touchcancel', this.stopInteraction);
  },
  beforeDestroy() {
    window.removeEventListener('resize', this.detectMobile);
    window.removeEventListener('mousemove', this.throttledMouseMove);
    window.removeEventListener('mouseup', this.stopInteraction);

    window.removeEventListener('touchmove', this.onTouchMove);
    window.removeEventListener('touchend', this.stopInteraction);
    window.removeEventListener('touchcancel', this.stopInteraction);
  },
  methods: {
    findNode(id) {
      return this.nodes.find(node => node.id === id);
    },
    selectNode(id) {
      if (this.mode === 'remove') {
        this.removeNode(id);
      } else if (this.mode === 'edit') {
        if (id === 'S' || id === 'P') {
          alert('Node cannot be updated.');
          this.edgeStartNode = null;
          this.selectedNodeId = null;
          this.hoveredNodeId = null;
          this.selectedNodeIdInEditMode = null;
          return;
        }

        this.selectedNodeIdInEditMode = id;
        const newId = prompt('Enter new ID', id);
        if (newId && newId !== id) {
          if (!this.findNode(newId)) {
            this.updateNodeId(id, newId);
          } else {
            alert('Duplicate ID detected. Please enter a unique ID.');
          }
        }
        this.selectedNodeIdInEditMode = null;
      } else if (this.mode === 'addEdge') {
        if (!this.edgeStartNode) {
          this.edgeStartNode = id;
          this.selectedNodeId = id;
        } else {
          if (this.edgeStartNode !== id) {
            const existingEdge = this.edges.find(edge => 
              (edge.from === this.edgeStartNode && edge.to === id) ||
              (edge.from === id && edge.to === this.edgeStartNode)
            );
            if (!existingEdge) {
              const weight = prompt('Enter weight for new edge');
              if (weight) {
                const parsedWeight = parseInt(weight);
                if (Number.isInteger(parsedWeight) && parsedWeight > 0) {
                  this.addEdge(this.edgeStartNode, id, parsedWeight);
                  this.edgeStartNode = null;
                  this.selectedNodeId = null;
                  this.hoveredNodeId = null;
                  this.selectedNodeIdInEditMode = null;
                } else {
                  alert('Weight must be a positive integer.');
                  this.edgeStartNode = null;
                  this.selectedNodeId = null;
                  this.hoveredNodeId = null;
                  this.selectedNodeIdInEditMode = null;
                }
              }
            } else {
                alert('An edge already exists between these nodes.');
                this.edgeStartNode = null;
                this.selectedNodeId = null;
                this.hoveredNodeId = null;
                this.selectedNodeIdInEditMode = null;
            }
          } else {
              alert('Cannot connect a node to itself. Select a different node.');
              this.edgeStartNode = null;
              this.selectedNodeId = null;
              this.hoveredNodeId = null;
              this.selectedNodeIdInEditMode = null;
          }
        }
      } else if (this.mode === 'pan' || this.mode === 'drag') {
        this.selectedNodeId = id;
        this.edgeStartNode = null;
        this.hoveredNodeId = null;
        this.selectedNodeIdInEditMode = null;
      }
    },
    handleEdgeSelect(edge) {
      if (this.mode === 'remove') {
        this.removeEdge(edge.id);
      } else if (this.mode === 'edit') {
        this.selectedEdgeId = edge.id;
        const newWeight = prompt('Enter new weight', edge.weight);
        if (newWeight) {
          this.updateEdgeWeight(edge.id, parseInt(newWeight));
        }
        this.selectedEdgeId = null;
      }
    },
    startNodeDrag(event, nodeId) {
      if (this.mode === 'drag') {
        event.preventDefault();
        this.draggingNodeId = nodeId;
        const node = this.findNode(nodeId);
        if (node) {
          if (event.type === 'touchstart' && event.touches.length === 1) {
            console.log('Mobile touch start');
            this.offsetX = event.targetTouches[0].clientX - node.x - this.panX;
            this.offsetY = event.targetTouches[0].clientY - node.y - this.panY;
          } else {
            this.offsetX = event.clientX - node.x - this.panX;
            this.offsetY = event.clientY - node.y - this.panY;
          }
          console.log(event.type, this.offsetX, this.offsetY);
          if (event.type === 'touchend') {
            this.draggingNodeId = null;
          }
        }
      }
    },
    startInteraction(event) {
      if (this.mode === 'pan') {
        if (event.type === 'mousedown' || (event.type === 'touchstart' && event.touches.length === 1)) {
          this.isPanning = true;
          this.panStartX = event.type === 'mousedown' ? event.clientX : event.touches[0].clientX;
          this.panStartY = event.type === 'mousedown' ? event.clientY : event.touches[0].clientY;
        }
      }
      if (this.mode === 'drag' && this.draggingNodeId) {
        this.selectedNodeId = this.draggingNodeId;
      }
    },
    stopInteraction() {
      this.isPanning = false;
      this.draggingNodeId = null;
      this.selectedEdgeId = null;
      this.selectedNodeId = null;
      this.highlightedEdgeId = null;
      this.selectedNodeIdInEditMode = null;
      this.addNodeHovered = null;
    },
    onMouseMove(event) {
      if (this.draggingNodeId !== null && this.mode === 'drag') {
        const node = this.findNode(this.draggingNodeId);
        if (node) {
          node.x = event.clientX - this.offsetX - this.panX;
          node.y = event.clientY - this.offsetY - this.panY;
        }
      }
      if (this.isPanning) {
        const dx = event.clientX - this.panStartX;
        const dy = event.clientY - this.panStartY;
        this.panX += dx;
        this.panY += dy;
        this.panStartX = event.clientX;
        this.panStartY = event.clientY;
      }
    },
    onTouchMove(event) {
      if (this.draggingNodeId !== null && this.mode === 'drag') {
        const node = this.findNode(this.draggingNodeId);
        if (node) {
          const touch = event.targetTouches[0];
          node.x = touch.clientX - this.offsetX - this.panX;
          node.y = touch.clientY - this.offsetY - this.panY;
        }
      }else if (this.isPanning) {
        let clientX, clientY;
        if (event.type === 'touchmove') {
          clientX = event.touches[0].clientX;
          clientY = event.touches[0].clientY;
        }

        const dx = clientX - this.panStartX;
        const dy = clientY - this.panStartY;
        this.panStartX = clientX;
        this.panStartY = clientY;

        this.panX += dx;
        this.panY += dy;
      }
    },
    handleMouseOver(event) {
      const svgRect = this.$refs.svg.getBoundingClientRect();
      const x = event.clientX - svgRect.left - this.panX;
      const y = event.clientY - svgRect.top - this.panY;
      const hoveredNode = this.nodes.find(node => Math.abs(node.x - x) < 20 && Math.abs(node.y - y) < 20);

      if (this.mode === 'addEdge' || this.mode === 'edit' || this.mode === 'drag' || this.mode === 'remove') {
        this.hoveredNodeId = hoveredNode ? hoveredNode.id : null;
        if ((this.mode === 'edit' || this.mode === 'remove') && !hoveredNode) {
          const hoveredEdge = this.edges.find(edge => {
            const fromNode = this.findNode(edge.from);
            const toNode = this.findNode(edge.to);
            if (!fromNode || !toNode) return false;
            const dist = Math.abs(
              (toNode.y - fromNode.y) * x - (toNode.x - fromNode.x) * y + toNode.x * fromNode.y - toNode.y * fromNode.x
            ) /
              Math.sqrt(
                Math.pow(toNode.y - fromNode.y, 2) + Math.pow(toNode.x - fromNode.x, 2)
              );
            return dist < 20;
          });
          this.highlightedEdgeId = hoveredEdge ? hoveredEdge.id : null;
        } else {
          this.highlightedEdgeId = null;
        }
      } else if (this.mode === 'add') {
        this.addNodeHovered = { x, y };
      } else {
        this.hoveredNodeId = null;
        this.addNodeHovered = null;
        this.highlightedEdgeId = null;
      }
      if (this.draggingNodeId === null) {
        this.selectedNodeId = this.hoveredNodeId ? this.hoveredNodeId : null;
      }
    },
    handleSvgClick(event) {
      if (this.mode === 'add') {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const x = event.clientX - svgRect.left - this.panX;
        const y = event.clientY - svgRect.top - this.panY;
        const newId = prompt('Enter new ID');
        if (newId) {
          if (!this.findNode(newId)) {
            this.addNode(newId, x, y);
          } else {
            alert('Duplicate ID detected. Please enter a unique ID.');
            this.edgeStartNode = null;
            this.selectedNodeId = null;
            this.hoveredNodeId = null;
            this.selectedNodeIdInEditMode = null;
          }
        }
      } else {
        if(this.mode !== 'addEdge') {
          this.selectedNodeId = null;
          this.edgeStartNode = null;
        }
      }
      this.addNodeHovered = null;
    },
    addNode(id, x, y) {
      if (!this.findNode(id)) {
        this.nodes.push({ id, x, y });
      } else {
        alert('Duplicate ID detected. Please enter a unique ID.');
      }
    },
    removeNode(id) {
      if (id === 'S' || id === 'P') {
        alert('Node cannot be deleted.');
        this.edgeStartNode = null;
        this.selectedNodeId = null;
        this.hoveredNodeId = null;
        this.selectedNodeIdInEditMode = null;
        return;
      }

      this.nodes = this.nodes.filter(node => node.id !== id);
      this.edges = this.edges.filter(edge => edge.from !== id && edge.to !== id);
    },
    addEdge(from, to, weight) {
      if (!this.edges.find(edge => (edge.from === from && edge.to === to) || (edge.from === to && edge.to === from))) {
        this.edges.push({ id: `${from}-${to}`, from, to, weight });
      }
    },
    removeEdge(id) {
      this.edges = this.edges.filter(edge => edge.id !== id);
    },
    updateNodeId(oldId, newId) {
      const node = this.findNode(oldId);
      if (node) {
        node.id = newId;
        this.edges.forEach(edge => {
          if (edge.from === oldId) edge.from = newId;
          if (edge.to === oldId) edge.to = newId;
        });
      }
    },
    updateEdgeWeight(id, newWeight) {
      const edge = this.edges.find(edge => edge.id === id);
      if (edge) {
        edge.weight = newWeight;
      }
    },
    setMode(mode) {
      this.mode = mode;
      this.isPanning = false;
      this.draggingNodeId = null;
      this.edgeStartNode = null;
      this.selectedNodeId = null;
      this.hoveredNodeId = null;
      this.highlightedEdgeId = null;
      this.addNodeHovered = null;
    },
    getGraphRelations() {
      const nodes = this.nodes.map(node => node.id);
      const edges = this.edges.map(edge => ({
        from: edge.from,
        to: edge.to,
        weight: edge.weight
      }));

      const graphRelations = {
        nodes,
        edges
      };
      console.log(graphRelations);

      return graphRelations;
    },
    getGraphRelationsAndPositions() {
      const nodes = this.nodes.map(node => ({
        id: node.id,
        x: node.x,
        y: node.y
      }));

      const edges = this.edges.map(edge => ({
        from: edge.from,
        to: edge.to,
        weight: edge.weight
      }));

      const graphRelationsAndPositions = {
        nodes,
        edges
      };
      console.log(graphRelationsAndPositions);

      return graphRelationsAndPositions;
    },

    async toggleAnimation() {
      this.animationError = false;
      
      if (!this.animationRunning) {
        this.animationRunning = true;
        await this.startAnimation();
      } else {
        this.animationRunning = false;
      }
    },

    async startAnimation() {
      try {
        const response = await fetch('http://localhost:8000/api/v1/graphs/shortest-path', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            nodes: this.nodes.map(node => node.id),
            edges: this.edges.map(edge => ({
              from: edge.from,
              to: edge.to,
              weight: edge.weight
            }))
          })
        });
        
        const data = await response.json();
        if (data.path) {
          this.animatedPath = [];
          this.animatedPath.push(this.findNode('S'));

          this.animatedPath.push(...data.path.map(nodeId => this.findNode(nodeId)));
          this.pathIndex = 0;
          
          this.animationRunning = true;
          this.animationError = false;

          this.animateToucan();
        } else {
          console.error('Path data is not valid:', data);
          setTimeout(() => { this.animationRunning = false; this.animationError = true; }, 200);
        }
      } catch (error) {
        console.error('Error fetching path:', error);
        setTimeout(() => { this.animationRunning = false; this.animationError = true; }, 200);
      }
    },

    async animateToucan() {
      if (this.animationRunning && this.pathIndex < this.animatedPath.length) {
        const node = this.animatedPath[this.pathIndex];
        await this.panToNode(node.id);
        this.pathIndex++;
        setTimeout(this.animateToucan, 1000);
      } else {
        this.animationRunning = false;
      }
    },

    async panToNode(nodeId) {
      const node = this.findNode(nodeId);
      if (node) {
        const svgRect = this.$refs.svg.getBoundingClientRect();
        const centerX = svgRect.width / 2;
        const centerY = svgRect.height / 2;

        const currentPanX = this.panX;
        const currentPanY = this.panY;

        const targetX = -node.x + centerX - currentPanX;
        const targetY = -node.y + centerY - currentPanY;

        const duration = 500;
        let startTime = null;

        const animate = (currentTime) => {
          if (!startTime) startTime = currentTime;
          const elapsedTime = currentTime - startTime;
          const progress = Math.min(elapsedTime / duration, 1);

          const easingProgress = this.easeInOutQuad(progress);

          this.panX = currentPanX + targetX * easingProgress;
          this.panY = currentPanY + targetY * easingProgress;

          const dx = (node.x + this.panX) - this.toucanX;
          const dy = (node.y + this.panY) - this.toucanY;
          this.toucanX += dx * easingProgress;
          this.toucanY += dy * easingProgress;

          if (progress < 1) {
            requestAnimationFrame(animate);
          }
        };

        requestAnimationFrame(animate);
      }
    },

    adjustBottomButtonGroup() {
      // Ensure bottom button group remains visible when the keyboard is open
      if (this.isMobile) {
        const bottomButtonGroup = document.querySelector('.button-group.bottom');
        if (bottomButtonGroup) {
          bottomButtonGroup.style.visibility = 'visible';
          setTimeout(() => {
            const rect = bottomButtonGroup.getBoundingClientRect();
            const isVisible = rect.bottom <= window.innerHeight;
            if (!isVisible) {
              bottomButtonGroup.style.visibility = 'hidden';
            }
          }, 300); // Adjust delay as needed for keyboard to open
        }
      }
    },

    detectMobile() {
      this.isMobile = window.innerWidth <= 768 || window.innerHeight <= 520;
    },

    easeInOutQuad(t) {
      return t < 0.5 ? 2 * t * t : -1 + (4 - 2 * t) * t;
    },

    handleTouchMove(event) {
      event.preventDefault();
      this.handleMouseOver(event.touches[0]);
    },
  }
};
</script>

<style scoped>
.fill-height {
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
}

.svg-container {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.graph-svg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1; /* Ensure SVG content is below buttons */
}

.button-group {
  position: absolute;
  z-index: 10;
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 8px;
}
.button-group.top {
  top: 16px;
  left: 16px;
}
.button-group.bottom {
  bottom: 16px;
  left: 16px;
}
.button-group.mobile {
  flex-direction: row;
  gap: 4px;
  padding: 4px;
}
.button-group.mobile .v-btn {
  min-width: auto; /* Adjust as needed */
  padding: 6px 12px; /* Adjust padding */
}
</style>
